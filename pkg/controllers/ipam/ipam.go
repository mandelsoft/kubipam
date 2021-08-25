/*
 * Copyright 2020 Mandelsoft. All rights reserved.
 *  This file is licensed under the Apache Software License, v. 2 except as noted
 *  otherwise in the LICENSE file
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package controllers

import (
	"fmt"
	"sync"

	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller/reconcile"
	"github.com/gardener/controller-manager-library/pkg/logger"
	"github.com/gardener/controller-manager-library/pkg/resources"

	api "github.com/mandelsoft/kubipam/pkg/apis/ipam/v1alpha1"
	"github.com/mandelsoft/kubipam/pkg/ipam"
)

type IPAM struct {
	lock      sync.RWMutex
	object    resources.Object
	ipam      *ipam.IPAM
	chunksize int
	error     string
	deleted   bool
}

func (this *Reconciler) setupIPAM(logger logger.LogContext, obj resources.Object) (bool, error) {
	r := obj.Data().(*api.IPAMRange)

	o := &IPAM{object: obj, chunksize: r.Spec.ChunkSize}
	this.ipams[obj.ObjectName()] = o

	ranges, err := ipam.ParseIPRanges(r.Spec.Ranges...)
	if err != nil {
		o.error = err.Error()
		return true, err
	}
	ipr, err := ipam.NewIPAMForRanges(ranges)
	if err != nil {
		o.error = err.Error()
		return true, err
	}
	if r.Spec.Mode == api.MODE_ROUNDROBIN {
		ipr.SetRoundRobin(true)
		ipr.SetState(nil, r.GetState())
	} else {
		ipr.SetRoundRobin(false)
	}
	o.ipam = ipr
	return true, nil
}

func (this *Reconciler) getRange(name resources.ObjectName) *IPAM {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.ipams[name]
}

func (this *Reconciler) setRange(name resources.ObjectName, ipr *IPAM) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.ipams[name] = ipr
}

func (this *Reconciler) reconcileRange(logger logger.LogContext, obj resources.Object) reconcile.Status {
	old := this.getRange(obj.ObjectName())

	if old != nil {
		logger.Infof("reconcile existing")
	} else {
		logger.Infof("reconcile new")
	}
	r := obj.Data().(*api.IPAMRange)
	ranges, err := ipam.ParseIPRanges(r.Spec.Ranges...)

	roundRobin := false
	if err == nil {
		switch r.Spec.Mode {
		case "", api.MODE_FIRSTMATCH:
			roundRobin = false
		case api.MODE_ROUNDROBIN:
			roundRobin = true
		default:
			err = fmt.Errorf("invalid mode %q: use %s or %s", r.Spec.Mode, api.MODE_FIRSTMATCH, api.MODE_ROUNDROBIN)
		}
	}

	var ipr *ipam.IPAM
	if err == nil {
		ipr, err = ipam.NewIPAMForRanges(ranges)
	}

	if err == nil {
		if ipr.Bits() < r.Spec.ChunkSize {
			err = fmt.Errorf("chunk size %d too large: network %d", r.Spec.ChunkSize, ipr.Bits())
		}
	}

	if err != nil {
		if old != nil {
			old.error = err.Error()
		}
		return reconcile.UpdateStatus(logger, resources.NewStandardStatusUpdate(logger, obj, api.STATE_INVALID, err.Error()))
	}

	if old == nil {
		new := &IPAM{
			object:    obj,
			ipam:      ipr,
			chunksize: r.Spec.ChunkSize,
			error:     "",
		}
		new.lock.Lock()
		defer new.lock.Unlock()
		this.setRange(obj.ObjectName(), new)
		new.ipam.SetRoundRobin(roundRobin)
	} else {
		old.lock.Lock()
		defer old.lock.Unlock()
		old.object = obj
		old.chunksize = r.Spec.ChunkSize
		old.ipam.SetRoundRobin(roundRobin)
	}
	if len(this.GetUsersFor(obj.ClusterKey())) > 0 {
		if !this.Controller().HasFinalizer(obj) {
			logger.Infof("setting finalizer because of pending requests")
			if err := this.Controller().SetFinalizer(obj); err != nil {
				return reconcile.Delay(logger, err)
			}
		}
	}
	if r.Spec.Mode == "" {
		mode := api.MODE_FIRSTMATCH
		if ipr.IsRoundRobin() {
			mode = api.MODE_ROUNDROBIN
		}
		reconcile.Update(logger, resources.NewUpdater(obj, func(mod *resources.ModificationState) error {
			r := mod.Data().(*api.IPAMRange)
			mod.AssureStringValue(&r.Spec.Mode, mode)
			return nil
		}))
	}
	return reconcile.UpdateStatus(logger, resources.NewStandardStatusUpdate(logger, obj, api.STATE_READY, ""))
}

func (this *Reconciler) deleteRange(logger logger.LogContext, obj resources.Object) reconcile.Status {
	old := this.getRange(obj.ObjectName())
	if old != nil {
		old.lock.Lock()
		defer old.lock.Unlock()
		if len(this.GetUsersFor(obj.ClusterKey())) > 0 {
			logger.Infof("mark as deleted")
			old.deleted = true
			old.error = "IPRange deleted"
			return reconcile.UpdateStatus(logger, resources.NewStandardStatusUpdate(logger, obj, api.STATE_DELETING,
				"waiting for pending requests to be deleted"))
		} else {
			if this.Controller().HasFinalizer(obj) {
				logger.Infof("removing finalizer because of no more requests")
				if err := this.Controller().RemoveFinalizer(obj); err != nil {
					return reconcile.Delay(logger, err)
				}
			}
		}
	}
	return reconcile.Succeeded(logger)
}

func (this *Reconciler) deletedRange(logger logger.LogContext, key resources.ClusterObjectKey) reconcile.Status {
	this.lock.Lock()
	defer this.lock.Unlock()
	logger.Infof("finally delete state")
	delete(this.ipams, key.ObjectName())
	return reconcile.Succeeded(logger)
}
