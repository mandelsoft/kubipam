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

package resources

import (
	"reflect"

	"github.com/gardener/controller-manager-library/pkg/resources/abstract"
	"github.com/gardener/controller-manager-library/pkg/resources/errors"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type AbstractResource struct {
	*abstract.AbstractResource
	helper *ResourceHelper
}

type ResourceHelper struct {
	Internal
}

func NewAbstractResource(ctx ResourceContext, self Internal, otype, ltype reflect.Type, gvk schema.GroupVersionKind) (AbstractResource, *ResourceHelper) {
	abs := abstract.NewAbstractResource(ctx, otype, ltype, gvk)
	helper := &ResourceHelper{self}
	return AbstractResource{abs, helper}, helper
}

func (this *AbstractResource) Name() string {
	return this.helper.Internal.Info().Name()
}

func (this *AbstractResource) Namespaced() bool {
	return this.helper.Internal.Info().Namespaced()
}

func (this *AbstractResource) Wrap(obj ObjectData) (Object, error) {
	if err := this.CheckOType(obj); err != nil {
		return nil, err
	}
	return this.helper.ObjectAsResource(obj), nil
}

func (this *AbstractResource) New(name ObjectName) Object {
	data := this.CreateData()
	data.GetObjectKind().SetGroupVersionKind(this.GroupVersionKind())
	if name != nil {
		data.SetName(name.Name())
		data.SetNamespace(name.Namespace())
	}
	return this.helper.ObjectAsResource(data)
}

func (this *AbstractResource) Namespace(namespace string) Namespaced {
	return &namespacedResource{this, namespace, nil}
}

////////////////////////////////////////////////////////////////////////////////
// Resource Helper

func (this *ResourceHelper) ObjectAsResource(obj ObjectData) Object {
	return newObject(obj, this.Internal)
}

func (this *ResourceHelper) Get(namespace, name string, result ObjectData) (Object, error) {
	if !this.Namespaced() && namespace != "" {
		return nil, errors.ErrNotNamespaced.New(this.GroupKind())
	}
	if this.Namespaced() && namespace == "" {
		return nil, errors.ErrNamespaced.New(this.GroupKind())
	}

	if result == nil {
		result = this.I_CreateData()
	}
	result.SetNamespace(namespace)
	result.SetName(name)
	err := this.I_get(result)
	return this.ObjectAsResource(result), err
}
