package solvers

import (
	controllerv1alpha1 "github.com/devfile/devworkspace-operator/apis/controller/v1alpha1"
	"github.com/devfile/devworkspace-operator/controllers/controller/workspacerouting/solvers"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ExampleSolver struct{}

var _ solvers.RoutingSolver = (*ExampleSolver)(nil)

func (s *ExampleSolver) GetSpecObjects(_ *controllerv1alpha1.WorkspaceRouting, workspaceMeta solvers.WorkspaceMetadata) solvers.RoutingObjects {
	exampleService := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-service",
			Namespace: workspaceMeta.Namespace,
			Labels: map[string]string{
				"controller.devfile.io/workspace_id": workspaceMeta.WorkspaceId,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(9999),
					TargetPort: intstr.FromInt(9999),
				},
			},
			Selector: workspaceMeta.PodSelector,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	return solvers.RoutingObjects{
		Services:     []corev1.Service{exampleService},
		Ingresses:    nil,
		Routes:       nil,
		PodAdditions: nil,
		OAuthClient:  nil,
	}
}

func (s *ExampleSolver) GetExposedEndpoints(endpoints map[string]controllerv1alpha1.EndpointList, routingObj solvers.RoutingObjects) (exposedEndpoints map[string]controllerv1alpha1.ExposedEndpointList, ready bool, err error) {
	service := routingObj.Services[0]
	exposedEndpoints = map[string]controllerv1alpha1.ExposedEndpointList{}
	ready = true

	for machineName, machineEndpoints := range endpoints {
		for _, endpoint := range machineEndpoints {
			endpointAttributes := map[string]string{}
			err := endpoint.Attributes.Into(&endpointAttributes)
			if err != nil {
				return nil, false, err
			}
			exposedEndpoints[machineName] = append(exposedEndpoints[machineName], controllerv1alpha1.ExposedEndpoint{
				Name:       endpoint.Name,
				Url:        service.Spec.ClusterIP,
				Attributes: endpointAttributes,
			})
		}
	}
	return exposedEndpoints, true, nil
}

type ExampleRoutingGetter struct{}

var _ solvers.RoutingSolverGetter = (*ExampleRoutingGetter)(nil)

func (e ExampleRoutingGetter) HasSolver(routingClass controllerv1alpha1.WorkspaceRoutingClass) bool {
	switch routingClass {
	case "external-sample":
		return true
	default:
		return false
	}
}

func (e ExampleRoutingGetter) GetSolver(_ client.Client, routingClass controllerv1alpha1.WorkspaceRoutingClass) (solver solvers.RoutingSolver, err error) {
	switch routingClass {
	case "external-sample":
		return &ExampleSolver{}, nil
	default:
		return nil, solvers.RoutingNotSupported
	}
}
