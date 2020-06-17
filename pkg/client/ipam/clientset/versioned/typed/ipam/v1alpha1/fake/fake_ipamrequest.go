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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/mandelsoft/kubipam/pkg/apis/ipam/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeIPAMRequests implements IPAMRequestInterface
type FakeIPAMRequests struct {
	Fake *FakeIpamV1alpha1
	ns   string
}

var ipamrequestsResource = schema.GroupVersionResource{Group: "ipam.mandelsoft.org", Version: "v1alpha1", Resource: "ipamrequests"}

var ipamrequestsKind = schema.GroupVersionKind{Group: "ipam.mandelsoft.org", Version: "v1alpha1", Kind: "IPAMRequest"}

// Get takes name of the iPAMRequest, and returns the corresponding iPAMRequest object, and an error if there is any.
func (c *FakeIPAMRequests) Get(name string, options v1.GetOptions) (result *v1alpha1.IPAMRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(ipamrequestsResource, c.ns, name), &v1alpha1.IPAMRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IPAMRequest), err
}

// List takes label and field selectors, and returns the list of IPAMRequests that match those selectors.
func (c *FakeIPAMRequests) List(opts v1.ListOptions) (result *v1alpha1.IPAMRequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(ipamrequestsResource, ipamrequestsKind, c.ns, opts), &v1alpha1.IPAMRequestList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.IPAMRequestList{ListMeta: obj.(*v1alpha1.IPAMRequestList).ListMeta}
	for _, item := range obj.(*v1alpha1.IPAMRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested iPAMRequests.
func (c *FakeIPAMRequests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(ipamrequestsResource, c.ns, opts))

}

// Create takes the representation of a iPAMRequest and creates it.  Returns the server's representation of the iPAMRequest, and an error, if there is any.
func (c *FakeIPAMRequests) Create(iPAMRequest *v1alpha1.IPAMRequest) (result *v1alpha1.IPAMRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(ipamrequestsResource, c.ns, iPAMRequest), &v1alpha1.IPAMRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IPAMRequest), err
}

// Update takes the representation of a iPAMRequest and updates it. Returns the server's representation of the iPAMRequest, and an error, if there is any.
func (c *FakeIPAMRequests) Update(iPAMRequest *v1alpha1.IPAMRequest) (result *v1alpha1.IPAMRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(ipamrequestsResource, c.ns, iPAMRequest), &v1alpha1.IPAMRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IPAMRequest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeIPAMRequests) UpdateStatus(iPAMRequest *v1alpha1.IPAMRequest) (*v1alpha1.IPAMRequest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(ipamrequestsResource, "status", c.ns, iPAMRequest), &v1alpha1.IPAMRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IPAMRequest), err
}

// Delete takes name of the iPAMRequest and deletes it. Returns an error if one occurs.
func (c *FakeIPAMRequests) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(ipamrequestsResource, c.ns, name), &v1alpha1.IPAMRequest{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeIPAMRequests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(ipamrequestsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.IPAMRequestList{})
	return err
}

// Patch applies the patch and returns the patched iPAMRequest.
func (c *FakeIPAMRequests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.IPAMRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(ipamrequestsResource, c.ns, name, pt, data, subresources...), &v1alpha1.IPAMRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IPAMRequest), err
}
