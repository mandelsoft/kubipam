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
	"sync"

	"github.com/gardener/controller-manager-library/pkg/config"
	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller/reconcile"
	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller/reconcile/reconcilers"
	"github.com/gardener/controller-manager-library/pkg/logger"
	"github.com/gardener/controller-manager-library/pkg/resources"

	api "github.com/mandelsoft/kubipam/pkg/apis/ipam/v1alpha1"
)

type Reconciler struct {
	reconcilers.ReconcilerSupport
	config *Config
	*reconcilers.SimpleUsageCache

	lock  sync.RWMutex
	ipams map[resources.ObjectName]*IPAM
}

var _ reconcile.Interface = &Reconciler{}

///////////////////////////////////////////////////////////////////////////////

func (this *Reconciler) Config() config.OptionSource {
	return this.config
}

///////////////////////////////////////////////////////////////////////////////

func (this *Reconciler) Setup() {
	resc, _ := this.Controller().GetMainCluster().Resources().Get(api.IPAMRANGE)
	reconcilers.ProcessResource(this.Controller(), "setup", resc, this.setupIPAM)
	resc, _ = this.Controller().GetMainCluster().Resources().Get(api.IPAMREQUEST)
	this.SimpleUsageCache.SetupFor(this.Controller(), resc, this.setupRequest)
	this.Controller().Infof("setup done")
}

func (this *Reconciler) Reconcile(logger logger.LogContext, obj resources.Object) reconcile.Status {
	switch obj.GroupKind() {
	case api.IPAMREQUEST:
		return this.reconcileRequest(logger, obj)
	case api.IPAMRANGE:
		return this.reconcileRange(logger, obj)
	}
	return reconcile.Succeeded(logger)
}

func (this *Reconciler) Delete(logger logger.LogContext, obj resources.Object) reconcile.Status {
	switch obj.GroupKind() {
	case api.IPAMREQUEST:
		return this.deleteRequest(logger, obj)
	case api.IPAMRANGE:
		return this.deleteRange(logger, obj)
	}
	return reconcile.Succeeded(logger)
}

func (this *Reconciler) Deleted(logger logger.LogContext, key resources.ClusterObjectKey) reconcile.Status {
	switch key.GroupKind() {
	case api.IPAMRANGE:
		return this.deletedRange(logger, key)
	case api.IPAMREQUEST:
		return this.deletedRequest(logger, key)
	}
	return reconcile.Succeeded(logger)
}
