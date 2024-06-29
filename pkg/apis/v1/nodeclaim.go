/*
Copyright The Kubernetes Authors.

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

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// NodeClaimSpec describes the desired state of the NodeClaim
type NodeClaimSpec struct {
	// Taints will be applied to the NodeClaim's node.
	// +optional
	Taints []v1.Taint `json:"taints,omitempty"`
	// StartupTaints are taints that are applied to nodes upon startup which are expected to be removed automatically
	// within a short period of time, typically by a DaemonSet that tolerates the taint. These are commonly used by
	// daemonsets to allow initialization and enforce startup ordering.  StartupTaints are ignored for provisioning
	// purposes in that pods are not required to tolerate a StartupTaint in order to have nodes provisioned for them.
	// +optional
	StartupTaints []v1.Taint `json:"startupTaints,omitempty"`
	// Requirements are layered with GetLabels and applied to every node.
	// +kubebuilder:validation:XValidation:message="requirements with operator 'In' must have a value defined",rule="self.all(x, x.operator == 'In' ? x.values.size() != 0 : true)"
	// +kubebuilder:validation:XValidation:message="requirements operator 'Gt' or 'Lt' must have a single positive integer value",rule="self.all(x, (x.operator == 'Gt' || x.operator == 'Lt') ? (x.values.size() == 1 && int(x.values[0]) >= 0) : true)"
	// +kubebuilder:validation:XValidation:message="requirements with 'minValues' must have at least that many values specified in the 'values' field",rule="self.all(x, (x.operator == 'In' && has(x.minValues)) ? x.values.size() >= x.minValues : true)"
	// +kubebuilder:validation:MaxItems:=100
	// +required
	Requirements []NodeSelectorRequirementWithMinValues `json:"requirements" hash:"ignore"`
	// Resources models the resource requirements for the NodeClaim to launch
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty" hash:"ignore"`
	// NodeClassRef is a reference to an object that defines provider specific configuration
	// +required
	NodeClassRef *NodeClassReference `json:"nodeClassRef"`
}

// A node selector requirement with min values is a selector that contains values, a key, an operator that relates the key and values
// and minValues that represent the requirement to have at least that many values.
type NodeSelectorRequirementWithMinValues struct {
	v1.NodeSelectorRequirement `json:",inline"`
	// This field is ALPHA and can be dropped or replaced at any time
	// MinValues is the minimum number of unique values required to define the flexibility of the specific requirement.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=50
	// +optional
	MinValues *int `json:"minValues,omitempty"`
}

// ResourceRequirements models the required resources for the NodeClaim to launch
// Ths will eventually be transformed into v1.ResourceRequirements when we support resources.limits
type ResourceRequirements struct {
	// Requests describes the minimum required resources for the NodeClaim to launch
	// +optional
	Requests v1.ResourceList `json:"requests,omitempty"`
}

type NodeClassReference struct {
	// Kind of the referent; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds"
	// +required
	Kind string `json:"kind"`
	// Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
	// +required
	Name string `json:"name"`
	// API version of the referent
	// +required
	Group string `json:"group"`
}

// +kubebuilder:object:generate=false
type Provider = runtime.RawExtension

// NodeClaim is the Schema for the NodeClaims API
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=nodeclaims,scope=Cluster,categories=karpenter
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".metadata.labels.node\\.kubernetes\\.io/instance-type",description=""
// +kubebuilder:printcolumn:name="Capacity",type="string",JSONPath=".metadata.labels.karpenter\\.sh/capacity-type",description=""
// +kubebuilder:printcolumn:name="Zone",type="string",JSONPath=".metadata.labels.topology\\.kubernetes\\.io/zone",description=""
// +kubebuilder:printcolumn:name="Node",type="string",JSONPath=".status.nodeName",description=""
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description=""
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description=""
// +kubebuilder:printcolumn:name="ID",type="string",JSONPath=".status.providerID",priority=1,description=""
// +kubebuilder:printcolumn:name="NodePool",type="string",JSONPath=".metadata.labels.karpenter\\.sh/nodepool",priority=1,description=""
// +kubebuilder:printcolumn:name="NodeClass",type="string",JSONPath=".spec.nodeClassRef.name",priority=1,description=""
type NodeClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="spec is immutable"
	// +required
	Spec   NodeClaimSpec   `json:"spec"`
	Status NodeClaimStatus `json:"status,omitempty"`
}

// NodeClaimList contains a list of NodeClaims
// +kubebuilder:object:root=true
type NodeClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeClaim `json:"items"`
}
