package context

import (
	"context"
	"github.com/redhat-developer/service-binding-operator/apis/spec/v1alpha2"
	"github.com/redhat-developer/service-binding-operator/pkg/reconcile/pipeline"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var _ pipeline.Context = &specImpl{}

type specImpl struct {
	impl
	serviceBinding *v1alpha2.ServiceBinding
}

func (i *specImpl) BindingName() string {
	if i.serviceBinding.Spec.Name != "" {
		return i.serviceBinding.Spec.Name
	}
	return i.bindingMeta.Name
}

func (i *specImpl) Services() ([]pipeline.Service, error) {
	if i.services == nil {
		serviceRef := i.serviceBinding.Spec.Service.AsRefferable()

		gvr, err := i.typeLookup.ResourceForReferable(serviceRef)
		if err != nil {
				return nil, err
			}
		if serviceRef.Namespace == nil {
			serviceRef.Namespace = &i.serviceBinding.Namespace
		}
		u, err := i.client.Resource(*gvr).Namespace(*serviceRef.Namespace).Get(context.Background(), serviceRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		i.services = append(i.services, &service{client: i.client, resource: u, groupVersionResource: gvr, serviceRef: serviceRef})
	}
	services := make([]pipeline.Service, len(i.services))
	for idx := 0; idx < len(i.services); idx++ {
		services[idx] = i.services[idx]
	}
	return services, nil

}

func (i *specImpl) Applications() ([]pipeline.Application, error) {
	if i.applications == nil {
		ref := i.serviceBinding.Spec.Application.AsRefferable()
		gvr, err := i.typeLookup.ResourceForReferable(ref)
		if err != nil {
			return nil, err
		}
		if i.serviceBinding.Spec.Application.Name != "" {
			u, err := i.client.Resource(*gvr).Namespace(i.serviceBinding.Namespace).Get(context.Background(), ref.Name, metav1.GetOptions{})
			if err != nil {
				return nil, emptyApplicationsErr{err}
			}
			i.applications = append(i.applications, &application{gvr: gvr, persistedResource: u})
		}
		if i.serviceBinding.Spec.Application.Selector.MatchLabels != nil {
			matchLabels := i.serviceBinding.Spec.Application.Selector.MatchLabels
			opts := metav1.ListOptions{
				LabelSelector: labels.Set(matchLabels).String(),
			}

			objList, err := i.client.Resource(*gvr).Namespace(i.serviceBinding.Namespace).List(context.Background(), opts)
			if err != nil {
				return nil, err
			}

			if len(objList.Items) == 0 {
				return nil, emptyApplicationsErr{}
			}

			for index := range objList.Items {
				i.applications = append(i.applications, &application{gvr: gvr, persistedResource: &(objList.Items[index])})
			}
		}
	}

	result := make([]pipeline.Application, len(i.applications))
	for l, a := range i.applications {
		result[l] = a
	}
	return result, nil

}

func (s *specImpl) BindAsFiles() bool {
	return true
}

func (s *specImpl) MountPath() string {
	return ""
}

func (s *specImpl) NamingTemplate() string {
	return ""
}

func (s *specImpl) Mappings() map[string]string {
	return make(map[string]string)
}
