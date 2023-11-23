/*
Copyright 2022.

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

package v1beta1

import (
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// GlanceDefaults -
type GlanceDefaults struct {
	ContainerImageURL string
}

var glanceDefaults GlanceDefaults

// log is for logging in this package.
var glancelog = logf.Log.WithName("glance-resource")

// SetupGlanceDefaults - initialize Glance spec defaults for use with either internal or external webhooks
func SetupGlanceDefaults(defaults GlanceDefaults) {
	glanceDefaults = defaults
	glancelog.Info("Glance defaults initialized", "defaults", defaults)
}

// SetupWebhookWithManager sets up the webhook with the Manager
func (r *Glance) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-glance-openstack-org-v1beta1-glance,mutating=true,failurePolicy=fail,sideEffects=None,groups=glance.openstack.org,resources=glances,verbs=create;update,versions=v1beta1,name=mglance.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Glance{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Glance) Default() {
	glancelog.Info("default", "name", r.Name)

	r.Spec.Default()
}

// Default - set defaults for this Glance spec
func (spec *GlanceSpec) Default() {
	if len(spec.ContainerImage) == 0 {
		spec.ContainerImage = glanceDefaults.ContainerImageURL
	}
	if spec.GlanceAPIs == nil || len(spec.GlanceAPIs) == 0 {
		spec.GlanceAPIs = map[string]GlanceAPITemplate{}
	}
	for key, glanceAPI := range spec.GlanceAPIs {
		// Check the sub-cr ContainerImage parameter
		if glanceAPI.ContainerImage == "" {
			glanceAPI.ContainerImage = glanceDefaults.ContainerImageURL
			spec.GlanceAPIs[key] =  glanceAPI
		}
	}
}

//+kubebuilder:webhook:path=/validate-glance-openstack-org-v1beta1-glance,mutating=false,failurePolicy=fail,sideEffects=None,groups=glance.openstack.org,resources=glances,verbs=create;update,versions=v1beta1,name=vglance.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Glance{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Glance) ValidateCreate() error {
	glancelog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Glance) ValidateUpdate(old runtime.Object) error {
	glancelog.Info("validate update", "name", r.Name)

	// Type can either be "split" or "single": we do not support changing layout
	// because there's no logic in the operator to scale down the existing statefulset
	// and scale up the new one, hence updating the Spec.GlanceAPI.Type is not supported
	o := old.(*Glance);
	for key, glanceAPI := range r.Spec.GlanceAPIs {
		if (glanceAPI.Type != o.Spec.GlanceAPIs[key].Type) {
			return errors.New("GlanceAPI deployment layout can't be updated")
		}
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Glance) ValidateDelete() error {
	glancelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
