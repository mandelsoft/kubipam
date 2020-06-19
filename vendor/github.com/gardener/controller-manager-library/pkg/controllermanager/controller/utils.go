/*
 * Copyright 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 *
 */

package controller

import (
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/gardener/controller-manager-library/pkg/resources"
)

func (this *controller) HasFinalizer(obj resources.Object) bool {
	return this.finalizer.HasFinalizer(obj)
}

func (this *controller) SetFinalizer(obj resources.Object) error {
	return this.finalizer.SetFinalizer(obj)
}

func (this *controller) RemoveFinalizer(obj resources.Object) error {
	return this.finalizer.RemoveFinalizer(obj)
}

func (this *controller) FinalizerHandler() Finalizer {
	return this.finalizer
}
func (this *controller) SetFinalizerHandler(f Finalizer) {
	this.finalizer = f
}

///////////////////////////////////////////////////////////////////////////////

type definition_field interface {
	Definition
}

type DefinitionWrapper struct {
	definition_field
	filters []ResourceFilter
}

func (this *DefinitionWrapper) Definition() Definition {
	return this
}

func (this *DefinitionWrapper) ResourceFilters() []ResourceFilter {
	return append(this.ResourceFilters(), this.filters...)
}

var _ Definition = &DefinitionWrapper{}

func AddFilters(def Definition, filters ...ResourceFilter) Definition {
	return &DefinitionWrapper{def, filters}
}

func FinalizerName(domain, controller string) string {
	if domain == "" {
		return "acme.com" + "/" + controller
	}
	return domain + "/" + controller
}

////////////////////////////////////////////////////////////////////////////////

type WatchedResources map[string]resources.GroupKindSet

func (this WatchedResources) Add(key string, gks ...schema.GroupKind) WatchedResources {
	return this.GatheredAdd(key, nil, gks...)
}

func (this WatchedResources) Remove(key string, gks ...schema.GroupKind) WatchedResources {
	return this.GatheredRemove(key, nil, gks...)
}

func (this WatchedResources) Contains(key string, gk schema.GroupKind) bool {
	return this[key].Contains(gk)
}

func (this WatchedResources) GatheredAdd(key string, added resources.GroupKindSet, gks ...schema.GroupKind) WatchedResources {
	set := this[key]
	if set == nil {
		set = resources.GroupKindSet{}
		this[key] = set
	}
	for _, gk := range gks {
		if added != nil {
			if !set.Contains(gk) {
				added.Add(gk)
			} else {
				continue
			}
		}
		set.Add(gk)
	}
	return this
}

func (this WatchedResources) GatheredRemove(key string, removed resources.GroupKindSet, gks ...schema.GroupKind) WatchedResources {
	set := this[key]
	if set != nil {
		for _, gk := range gks {
			if removed != nil {
				if set.Contains(gk) {
					removed.Add(gk)
				} else {
					continue
				}
			}
			set.Remove(gk)
		}
		if len(set) == 0 {
			delete(this, key)
		}
	}
	return this
}
