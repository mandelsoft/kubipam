/*
Copyright (c) 2020 Mandelsoft. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/mandelsoft/kubipam/pkg/apis/ipam/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// IPAMRequestLister helps list IPAMRequests.
type IPAMRequestLister interface {
	// List lists all IPAMRequests in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.IPAMRequest, err error)
	// IPAMRequests returns an object that can list and get IPAMRequests.
	IPAMRequests(namespace string) IPAMRequestNamespaceLister
	IPAMRequestListerExpansion
}

// iPAMRequestLister implements the IPAMRequestLister interface.
type iPAMRequestLister struct {
	indexer cache.Indexer
}

// NewIPAMRequestLister returns a new IPAMRequestLister.
func NewIPAMRequestLister(indexer cache.Indexer) IPAMRequestLister {
	return &iPAMRequestLister{indexer: indexer}
}

// List lists all IPAMRequests in the indexer.
func (s *iPAMRequestLister) List(selector labels.Selector) (ret []*v1alpha1.IPAMRequest, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.IPAMRequest))
	})
	return ret, err
}

// IPAMRequests returns an object that can list and get IPAMRequests.
func (s *iPAMRequestLister) IPAMRequests(namespace string) IPAMRequestNamespaceLister {
	return iPAMRequestNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// IPAMRequestNamespaceLister helps list and get IPAMRequests.
type IPAMRequestNamespaceLister interface {
	// List lists all IPAMRequests in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.IPAMRequest, err error)
	// Get retrieves the IPAMRequest from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.IPAMRequest, error)
	IPAMRequestNamespaceListerExpansion
}

// iPAMRequestNamespaceLister implements the IPAMRequestNamespaceLister
// interface.
type iPAMRequestNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all IPAMRequests in the indexer for a given namespace.
func (s iPAMRequestNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.IPAMRequest, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.IPAMRequest))
	})
	return ret, err
}

// Get retrieves the IPAMRequest from the indexer for a given namespace and name.
func (s iPAMRequestNamespaceLister) Get(name string) (*v1alpha1.IPAMRequest, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("ipamrequest"), name)
	}
	return obj.(*v1alpha1.IPAMRequest), nil
}
