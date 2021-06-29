// +kubebuilder:object:generate=true
package apis

import (
	"errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Object reference in the same namespace
// +mapType=atomic
type Ref struct {

	// Group of the referent.
	Group string `json:"group"`

	// Version of the referent.
	Version string `json:"version"`

	// Kind of the referent.
	// +optional
	Kind string `json:"kind,omitempty"`

	// Resource of the referent.
	// +optional
	Resource string `json:"resource,omitempty"`

	// Name of the referent.
	Name string `json:"name,omitempty"`
}

// Object reference in some namespace
type NamespacedRef struct {
	Ref `json:",inline"`

	// Namespace of the referent.
	// if empty assumes the same namespace as ServiceBinding
	// +optional
	Namespace *string `json:"namespace,omitempty"`
}

// Returns GVR of reference if available, otherwise error
func (ref *Ref) GroupVersionResource() (*schema.GroupVersionResource, error) {
	if ref.Resource == "" {
		return nil, errors.New("Resource undefined")
	}
	return &schema.GroupVersionResource{
		Group:    ref.Group,
		Version:  ref.Version,
		Resource: ref.Resource,
	}, nil
}

// Returns GVK of reference if available, otherwise error
func (ref *Ref) GroupVersionKind() (*schema.GroupVersionKind, error) {
	if ref.Kind == "" {
		return nil, errors.New("Kind undefined")
	}
	return &schema.GroupVersionKind{
		Group:   ref.Group,
		Version: ref.Version,
		Kind:    ref.Kind,
	}, nil
}
