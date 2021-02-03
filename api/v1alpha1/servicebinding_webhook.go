/*
Copyright 2021.

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

package v1alpha1

import (
	"context"
	"encoding/json"
	"github.com/go-logr/logr"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const requesterAnnotationKey = "service.binding/requester"

func (r *ServiceBinding) SetupWebhookWithManager(mgr ctrl.Manager) error {
	mgr.GetWebhookServer().Register("/mutate-operators-coreos-com-v1alpha1-servicebinding", &webhook.Admission{
		Handler: &admisionHandler{},
	})
	return nil
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-operators-coreos-com-v1alpha1-servicebinding,mutating=true,failurePolicy=fail,sideEffects=None,groups=operators.coreos.com,resources=servicebindings,verbs=create;update,versions=v1alpha1,name=mservicebinding.kb.io,admissionReviewVersions={v1beta1}

type admisionHandler struct {
	decoder *admission.Decoder
	log logr.Logger
}

var _ webhook.AdmissionHandler = &admisionHandler{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (ah *admisionHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	sb := &ServiceBinding{}
	err := ah.decoder.Decode(req, sb)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	sb.requester(req.UserInfo.Username)
	ah.log.Info("XXX", "requester", sb.Requester())
	marshaledSB , err:= json.Marshal(sb)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledSB)
}

func (ah *admisionHandler) InjectDecoder(decoder *admission.Decoder) error {
	ah.decoder = decoder
	return nil
}

func (ah *admisionHandler) InjectLogger(l logr.Logger) error {
	ah.log = l
	return nil
}

func (sb *ServiceBinding) requester(username string) {
	metav1.SetMetaDataAnnotation(&sb.ObjectMeta, requesterAnnotationKey, username)
}

// Return username of requester who submitted the service binding
func (sb *ServiceBinding) Requester() *string {
	req, found := sb.Annotations[requesterAnnotationKey]
	if found {
		return &req
	}
	return nil
}