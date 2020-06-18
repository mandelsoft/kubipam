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
	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller/reconcile"
	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller/reconcile/reconcilers"
	"github.com/gardener/controller-manager-library/pkg/resources"
	"github.com/gardener/controller-manager-library/pkg/resources/apiextensions"

	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller"

	"github.com/mandelsoft/kubipam/pkg/apis/ipam/crds"
	api "github.com/mandelsoft/kubipam/pkg/apis/ipam/v1alpha1"
)

const NAME = "ipam"

func init() {
	crds.AddToRegistry(apiextensions.DefaultRegistry())
}

func init() {
	controller.Configure(NAME).
		FinalizerDomain("mandelsoft.org").
		DefaultWorkerPool(5, 0).
		OptionsByExample("options", &Config{}).
		Reconciler(Create).
		MainResourceByGK(api.IPAMRANGE).
		WatchesByGK(api.IPAMREQUEST).
		With(reconcilers.UsageReconcilerForGKs("ipam", controller.CLUSTER_MAIN, api.IPAMRANGE)).
		MustRegister()
}

///////////////////////////////////////////////////////////////////////////////

func Create(controller controller.Interface) (reconcile.Interface, error) {
	cfg, err := controller.GetOptionSource("options")
	if err != nil {
		return nil, err
	}
	config := cfg.(*Config)

	return &Reconciler{
		ReconcilerSupport: reconcilers.NewReconcilerSupport(controller),
		config:            config,
		SimpleUsageCache:  reconcilers.GetSharedSimpleUsageCache(controller),
		ipams:             map[resources.ObjectName]*IPAM{},
	}, nil
}
