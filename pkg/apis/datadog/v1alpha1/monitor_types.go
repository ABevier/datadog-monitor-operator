package v1alpha1

import (
	"github.com/zorkian/go-datadog-api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MonitorSpec defines the desired state of Monitor
type MonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Query   string   `json:"query"`
	Message string   `json:"message"`
	Tags    []string `json:"tags,omitempty"`
}

// MonitorStatus defines the observed state of Monitor
type MonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ID int `json:"id"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Monitor is the Schema for the monitors API
// +kubebuilder:resource:path=monitors,scope=Namespaced
type Monitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MonitorSpec   `json:"spec,omitempty"`
	Status MonitorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MonitorList contains a list of Monitor
type MonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Monitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Monitor{}, &MonitorList{})
}

func (m *MonitorSpec) ToDDMonitor() *datadog.Monitor {
	tagCopy := make([]string, len(m.Tags))
	copy(tagCopy, m.Tags)

	return &datadog.Monitor{
		Type:    &m.Type,
		Query:   &m.Query,
		Name:    &m.Name,
		Message: &m.Message,
		Tags:    tagCopy,
	}
}
