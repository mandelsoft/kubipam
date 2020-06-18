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
	"net"
	"time"

	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller/reconcile"
	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller/reconcile/reconcilers"
	"github.com/gardener/controller-manager-library/pkg/fieldpath"
	"github.com/gardener/controller-manager-library/pkg/logger"
	"github.com/gardener/controller-manager-library/pkg/resources"
	corev1 "k8s.io/api/core/v1"

	api "github.com/mandelsoft/kubipam/pkg/apis/ipam/v1alpha1"
)

var assignedCIDRField = fieldpath.RequiredField(&api.IPAMRequest{}, ".Status.CIDR")
var rangeFilter = resources.NewGroupKindFilter(api.IPAMRANGE)

func (this *Reconciler) setupRequest(sub resources.Object) resources.ClusterObjectKeySet {
	req := sub.Data().(*api.IPAMRequest)
	ref := req.Spec.IPAM.RelativeTo(sub)
	if ref.Name() != "" {
		ipam := this.ipams[ref]
		if ipam != nil {
			if req.Status.CIDR != "" {
				_, cidr, err := net.ParseCIDR(req.Status.CIDR)
				if err != nil {
					this.Controller().Errorf("invalid state of ipam request %s: invalid cidr: %s", ref, req.Status.CIDR)
				} else {
					ipam.ipam.Busy(cidr)
				}
			}
		}
		return resources.NewClusterObjectKeySet(this.NewClusterObjectKey(api.IPAMRANGE, ref))
	}
	return nil
}

func (this *Reconciler) reconcileRequest(logger logger.LogContext, obj resources.Object) reconcile.Status {
	r := obj.Data().(*api.IPAMRequest)

	ref := r.Spec.IPAM.RelativeTo(obj)

	if ref.Name() == "" {
		return reconcile.UpdateStatus(logger, resources.NewStandardStatusUpdate(logger, obj, api.STATE_INVALID, "IPAMRange object not specified"))
	}

	if r.Spec.Request != "" {
		return reconcile.UpdateStatus(logger, resources.NewStandardStatusUpdate(logger, obj, api.STATE_INVALID, "request field not implemented yet"))
	}

	this.UpdateFilteredUsesFor(obj.ClusterKey(), rangeFilter, resources.NewClusterObjectKeySet(this.NewClusterObjectKey(api.IPAMRANGE, ref)))
	ipr := this.getRange(ref)
	if ipr == nil {
		return reconcile.UpdateStatus(logger, resources.NewStandardStatusUpdate(logger, obj, api.STATE_INVALID, fmt.Sprintf("IPAMRange %s not found", ref)))
	}
	if ipr.error != "" {
		return reconcile.UpdateStatus(logger, resources.NewStandardStatusUpdate(logger, obj, api.STATE_INVALID, fmt.Sprintf("IPAMRange %s not valid: %s", ref, ipr.error)))
	}

	ipr.lock.Lock()
	defer ipr.lock.Unlock()
	if r.Status.CIDR == "" {
		size := r.Spec.Size
		if size < 0 {
			return reconcile.UpdateStatus(logger, resources.NewStandardStatusUpdate(logger, obj, api.STATE_INVALID,
				fmt.Sprintf("invalid size %d", size)))
		}
		if size > ipr.ipam.Bits() {
			return reconcile.UpdateStatus(logger, resources.NewStandardStatusUpdate(logger, obj, api.STATE_INVALID,
				fmt.Sprintf("size %d too large: network %d", size, ipr.ipam.Bits())))
		}
		if size <= 0 {
			size = ipr.chunksize
		}
		if size <= 0 {
			size = ipr.ipam.Bits()
		}
		err := this.Controller().SetFinalizer(obj)
		if err != nil {
			return reconcile.Delay(logger, err)
		}
		if !this.Controller().HasFinalizer(ipr.object) {
			logger.Infof("requesting finalizer for IPAM %s", ref)
			if err := this.Controller().SetFinalizer(ipr.object); err != nil {
				return reconcile.Delay(logger, err)
			}
		}
		cidr := ipr.ipam.Alloc(size)
		if cidr != nil {
			logger.Infof("allocated %s", cidr)
			_, err := resources.ModifyStatus(obj, func(mod *resources.ModificationState) error {
				mod.Set(assignedCIDRField, cidr.String())
				return nil
			})
			if err != nil {
				ipr.ipam.Free(cidr)
				ipr.object.Event(corev1.EventTypeWarning, "allocation", fmt.Sprintf("allocation update failed: %s", err))
				return reconcile.Delay(logger, err)
			}
		} else {
			this.EnqueueKeys(this.GetUsesFor(this.NewClusterObjectKey(api.IPAMRANGE, ref)))
			ipr.object.Event(corev1.EventTypeWarning, "allocation", fmt.Sprintf("allocation size %d failed", size))
			return reconcile.UpdateStatus(logger, resources.NewStandardStatusUpdate(logger, obj, api.STATE_BUSY, fmt.Sprintf("requested chunk not available: %d", size)), 2*time.Minute)
		}
		ipr.object.Event(corev1.EventTypeNormal, "allocation", fmt.Sprintf("cidr %s allocated", cidr))
	}
	return reconcile.UpdateStatus(logger, resources.NewStandardStatusUpdate(logger, obj, api.STATE_READY, ""))
}

func (this *Reconciler) deleteRequest(logger logger.LogContext, obj resources.Object) reconcile.Status {
	if this.Controller().HasFinalizer(obj) {
		req := obj.Data().(*api.IPAMRequest)
		if req.Status.CIDR != "" {
			_, cidr, err := net.ParseCIDR(req.Status.CIDR)
			if err == nil {
				ref := req.Spec.IPAM.RelativeTo(obj)
				ipr := this.getRange(ref)
				if ipr != nil {
					ipr.lock.Lock()
					defer ipr.lock.Unlock()
					logger.Infof("releasing %s", cidr)
					ipr.ipam.Free(cidr)
					_, err := resources.Modify(obj, func(mod *resources.ModificationState) error {
						mod.Set(assignedCIDRField, "")
						return nil
					})
					if err != nil {
						ipr.ipam.Busy(cidr)
						ipr.object.Event(corev1.EventTypeWarning, "release", fmt.Sprintf("release update failed: %s", err))
						return reconcile.Delay(logger, err)
					}
					ipr.object.Event(corev1.EventTypeNormal, "release", fmt.Sprintf("cidr %s released", cidr))
				}
			}
		}
	}
	return reconcile.DelayOnError(logger, this.Controller().RemoveFinalizer(obj))
}

func (this *Reconciler) deletedRequest(logger logger.LogContext, key resources.ClusterObjectKey) reconcile.Status {
	this.CleanupUser(logger, "cleanup", this.Controller(), key, reconcilers.EnqueueAction)
	return reconcile.Succeeded(logger)
}
