/*
 * Copyright 2019 SAP SE or an SAP affiliate company. All rights reserved.
 * This file is licensed under the Apache Software License, v. 2 except as noted
 * otherwise in the LICENSE file
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

package plain

import (
	"github.com/gardener/controller-manager-library/pkg/resources"
	"github.com/gardener/controller-manager-library/pkg/resources/abstract"
)

type _object struct {
	*abstract.AbstractObject
}

var _ Object = &_object{}

func newObject(data ObjectData, resource Interface) Object {
	return &_object{abstract.NewAbstractObject(data, resource)}
}

func (this *_object) DeepCopy() Object {
	data := this.ObjectData.DeepCopyObject().(ObjectData)
	return newObject(data, this.GetResource())
}

/////////////////////////////////////////////////////////////////////////////////

func (this *_object) GetResource() Interface {
	return this.AbstractObject.GetResource().(Interface)
}

func (this *_object) Resources() Resources {
	return this.AbstractObject.GetResource().(Interface).Resources()
}

////////////////////////////////////////////////////////////////////////////////
// Modification

func (this *_object) ForCluster(cluster resources.Cluster) (resources.Object, error) {
	r, err := cluster.Resources().Get(this.GroupVersionKind())
	if err != nil {
		return nil, err
	}
	return r.Wrap(this.ObjectData)
}

func (this *_object) CreateIn(cluster resources.Cluster) error {
	o, err := this.ForCluster(cluster)
	if err != nil {
		return err
	}
	return o.Create()
}

func (this *_object) CreateOrUpdateIn(cluster resources.Cluster) error {
	o, err := this.ForCluster(cluster)
	if err != nil {
		return err
	}
	err = o.CreateOrUpdate()
	if err == nil {
		this.ObjectData = o.Data()
	}
	return err
}
func (this *_object) CreateOrModifyIn(cluster resources.Cluster, modifier resources.Modifier) (bool, error) {
	o, err := this.ForCluster(cluster)
	if err != nil {
		return false, err
	}
	mod, err := o.CreateOrModify(modifier)
	if err == nil {
		this.ObjectData = o.Data()
	}
	return mod, err
}

func (this *_object) UpdateIn(cluster resources.Cluster) error {
	o, err := this.ForCluster(cluster)
	if err != nil {
		return err
	}
	err = o.Update()
	if err == nil {
		this.ObjectData = o.Data()
	}
	return err
}

func (this *_object) ModifiyIn(cluster resources.Cluster, modifier resources.Modifier) (bool, error) {
	o, err := this.ForCluster(cluster)
	if err != nil {
		return false, err
	}
	mod, err := o.Modify(modifier)
	if err == nil {
		this.ObjectData = o.Data()
	}
	return mod, nil
}

func (this *_object) DeleteIn(cluster resources.Cluster) error {
	o, err := this.ForCluster(cluster)
	if err != nil {
		return err
	}
	return o.Delete()
}

func (this *_object) SetFinalizerIn(cluster resources.Cluster, key string) error {
	o, err := this.ForCluster(cluster)
	if err != nil {
		return err
	}
	err = o.SetFinalizer(key)
	if err == nil {
		this.ObjectData = o.Data()
	}
	return err
}
