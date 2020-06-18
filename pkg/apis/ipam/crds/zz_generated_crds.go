/*
Copyright (c) YEAR SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

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

package crds

import (
	"github.com/gardener/controller-manager-library/pkg/resources/apiextensions"
	"github.com/gardener/controller-manager-library/pkg/utils"
)

var registry = apiextensions.NewRegistry()

func init() {
	var data string
	data = `

---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.4
  creationTimestamp: null
  name: ipamranges.ipam.mandelsoft.org
spec:
  group: ipam.mandelsoft.org
  names:
    kind: IPAMRange
    listKind: IPAMRangeList
    plural: ipamranges
    shortNames:
    - iprange
    singular: ipamrange
  scope: Namespaced
  versions:
  - additionalPrinterColumns:
    - jsonPath: .status.state
      name: STATE
      type: string
    name: v1alpha1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            properties:
              chunkSize:
                type: integer
              ranges:
                items:
                  type: string
                type: array
            required:
            - ranges
            type: object
          status:
            properties:
              message:
                type: string
              state:
                type: string
            type: object
        required:
        - spec
        type: object
    served: true
    storage: true
    subresources:
      status: {}
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
  `
	utils.Must(registry.RegisterCRD(data))
	data = `

---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.4
  creationTimestamp: null
  name: ipamrequests.ipam.mandelsoft.org
spec:
  group: ipam.mandelsoft.org
  names:
    kind: IPAMRequest
    listKind: IPAMRequestList
    plural: ipamrequests
    shortNames:
    - ipreq
    singular: ipamrequest
  scope: Namespaced
  versions:
  - additionalPrinterColumns:
    - jsonPath: .spec.ipam.name
      name: IPAM
      type: string
    - jsonPath: .spec.size
      name: Size
      type: integer
    - jsonPath: .status.state
      name: STATE
      type: string
    - jsonPath: .status.cidr
      name: CIDR
      type: string
    name: v1alpha1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            properties:
              description:
                type: string
              ipam:
                description: ObjectReference is is plain reference to an object of
                  an implicitly determined type
                properties:
                  name:
                    type: string
                  namespace:
                    type: string
                required:
                - name
                type: object
              request:
                type: string
              size:
                type: integer
            required:
            - ipam
            type: object
          status:
            properties:
              cidr:
                type: string
              message:
                type: string
              state:
                type: string
            type: object
        required:
        - spec
        type: object
    served: true
    storage: true
    subresources:
      status: {}
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
  `
	utils.Must(registry.RegisterCRD(data))
}

func AddToRegistry(r apiextensions.Registry) {
	registry.AddToRegistry(r)
}
