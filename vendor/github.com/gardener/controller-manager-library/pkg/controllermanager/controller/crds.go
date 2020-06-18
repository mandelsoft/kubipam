/*
 * Copyright 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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
 *
 */

package controller

import (
	"github.com/Masterminds/semver"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/gardener/controller-manager-library/pkg/controllermanager/cluster"
	"github.com/gardener/controller-manager-library/pkg/resources/apiextensions"
	"github.com/gardener/controller-manager-library/pkg/utils"
)

type CustomResourceDefinition struct {
	versioned *utils.Versioned
}

func NewCustomResourceDefinition(crd ...*v1beta1.CustomResourceDefinition) *CustomResourceDefinition {
	if len(crd) > 1 {
		return nil
	}
	def := &CustomResourceDefinition{utils.NewVersioned(&v1beta1.CustomResourceDefinition{})}
	if len(crd) > 0 {
		def.versioned.SetDefault(crd[0])
	}
	return def
}

func (this *CustomResourceDefinition) GetFor(c cluster.Interface) *v1beta1.CustomResourceDefinition {
	f := this.versioned.GetFor(c.GetServerVersion())
	if f != nil {
		return f.(*v1beta1.CustomResourceDefinition)
	}
	return nil
}

func (this *CustomResourceDefinition) RegisterVersion(v *semver.Version, crd v1beta1.CustomResourceDefinition) *CustomResourceDefinition {
	this.versioned.MustRegisterVersion(v, crd)
	return this
}

func (this *CustomResourceDefinition) GroupKind() schema.GroupKind {
	if this.versioned.GetDefault() != nil {
		return this.versioned.GetDefault().(*v1beta1.CustomResourceDefinition).GroupVersionKind().GroupKind()
	}
	for _, o := range this.versioned.GetVersions() {
		return o.(*v1beta1.CustomResourceDefinition).GroupVersionKind().GroupKind()
	}
	return schema.GroupKind{}
}

// GetVersions provides an actual view for this deprecated type
func (this *CustomResourceDefinition) GetVersions() *apiextensions.CustomResourceDefinitionVersions {
	var vers *apiextensions.CustomResourceDefinitionVersions
	var err error
	if this.versioned.GetDefault() != nil {
		vers, err = apiextensions.NewDefaultedCustomResourceDefinitionVersions(this.versioned.GetDefault())
		utils.Must(err)
	} else {
		gk := this.GroupKind()
		if gk.Empty() {
			return nil
		}
		vers = apiextensions.NewCustomResourceDefinitionVersions(gk)
	}
	for v, o := range this.versioned.GetVersions() {
		utils.Must(vers.Override(v, o))
	}
	return vers
}
