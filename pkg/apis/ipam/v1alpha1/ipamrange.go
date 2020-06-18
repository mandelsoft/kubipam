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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type IPAMRangeList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IPAMRange `json:"items"`
}

// +kubebuilder:storageversion
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,path=ipamranges,shortName=iprange,singular=ipamrange
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name=CIDR,JSONPath=".spec.cidr",type=string
// +kubebuilder:printcolumn:name=ChunkSize,JSONPath=".spec.chunkSize",type=integer
// +kubebuilder:printcolumn:name=Busy,JSONPath=".status.busy",type=integer
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type IPAMRange struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              IPAMRangeSpec `json:"spec"`
	// +optional
	Status IPAMRangeStatus `json:"status,omitempty"`
}

type IPAMRangeSpec struct {
	CIDR      string   `json:"cidr"`
	ChunkSize []string `json:"chunkSize"`

	// +optional
	Ranges []string `json:"ranges"`
}

type IPAMRangeStatus struct {
	// +optional
	Busy string `json:"busy,omitempty"`
}
