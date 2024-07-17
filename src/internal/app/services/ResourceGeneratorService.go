package services

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kudejen/src/internal/api/server"
	resource_service "kudejen/src/internal/app/pb"
)

type ResourceGeneratorService struct {
	Server *server.Server
	resource_service.UnimplementedResourceGeneratorServer
}

func int32Ptr(i int32) *int32 { return &i }

var namespace string = "default"

func createPostgresDeployment(request *resource_service.CreateRequest) (*appsv1.Deployment, *v1.Service) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.GetName(),
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "postgres"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "postgres"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  request.GetName(),
							Image: "postgres:13",
							Ports: []v1.ContainerPort{
								{
									ContainerPort: 5432,
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "POSTGRES_DB",
									Value: request.GetDatabaseName(),
								},
								{
									Name:  "POSTGRES_USER",
									Value: request.GetUser(),
								},
								{
									Name:  "POSTGRES_PASSWORD",
									Value: request.GetPassword(),
								},
							},
						},
					},
				},
			},
		},
	}
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.GetName(),
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"app": "postgres"},
			Ports: []v1.ServicePort{
				{
					Port:       5432,
					TargetPort: intstr.FromInt(5432),
				},
			},
		},
	}
	return deployment, service
}
func (s ResourceGeneratorService) Create(ctx context.Context, request *resource_service.CreateRequest) (*resource_service.Response, error) {

	// Create TODO: use proto validator instead of manual request validation https://github.com/bufbuild/protovalidate
	if len(request.GetName()) < 3 || len(request.GetName()) > 50 ||
		len(request.GetUser()) < 3 || len(request.GetUser()) > 50 ||
		len(request.GetPassword()) < 3 || len(request.GetPassword()) > 50 ||
		len(request.GetDatabaseName()) < 3 || len(request.GetDatabaseName()) > 50 {
		return nil, fmt.Errorf("validation failure: svalues should be between 3 and 50 characters")
	}

	deployment, service := createPostgresDeployment(request)

	_, err := s.Server.KubernetesClient.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("service error: %s", err)
	}
	fmt.Println("Deployment created successfully")

	_, err = s.Server.KubernetesClient.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("service error: %s", err)
		//TODO:What happens if a deployment is created successfully, but the service creation fails
	}
	fmt.Println("Service created successfully")
	//TODO: maybe return some extra information from the deployment, we already have payload for that
	return &resource_service.Response{Message: "deployment has been created successfully!", Payload: make([]*resource_service.KeyValueItem, 0)}, nil
}

func (s ResourceGeneratorService) Update(ctx context.Context, request *resource_service.UpdateRequest) (*resource_service.Response, error) {
	// Create TODO: use proto validator instead of manual request validation https://github.com/bufbuild/protovalidate
	if len(request.GetName()) < 3 || len(request.GetName()) > 50 {
		return nil, fmt.Errorf("validation failure: service name should be between 3 and 50 characters")
	}
	if request.GetReplicas() < 1 || request.GetReplicas() > 3 {
		return nil, fmt.Errorf("validation failure: the number of replica should be between 1 and 3")
	}
	deployment, err := s.Server.KubernetesClient.AppsV1().Deployments(namespace).Get(context.TODO(), request.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("service error: %s", err)
	}

	// Modify the number of replicas
	deployment.Spec.Replicas = int32Ptr(request.GetReplicas()) // Set the desired number of replicas

	// Update the deployment
	updatedDeployment, err := s.Server.KubernetesClient.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("service error: %s", err)
	}

	fmt.Printf("Deployment replicas updated successfully: %s\n", updatedDeployment.Name)
	return &resource_service.Response{Message: fmt.Sprintf("Deployment replicas updated successfully: %s\n", updatedDeployment.Name), Payload: make([]*resource_service.KeyValueItem, 0)}, nil
}

func (s ResourceGeneratorService) Delete(ctx context.Context, request *resource_service.DeleteRequest) (*resource_service.Response, error) {
	// Create TODO: use proto validator instead of manual request validation https://github.com/bufbuild/protovalidate
	if len(request.GetName()) < 3 || len(request.GetName()) > 50 {
		return nil, fmt.Errorf("validation failure: service name should be between 3 and 50 characters")
	}

	err := s.Server.KubernetesClient.AppsV1().Deployments("default").Delete(context.TODO(), request.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("service error: %s", err)
	}
	err = s.Server.KubernetesClient.CoreV1().Services("default").Delete(context.TODO(), request.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("service error: %s", err)
	}
	return &resource_service.Response{Message: "the deployment and associated service have been deleted successfully", Payload: make([]*resource_service.KeyValueItem, 0)}, nil
}

func (s ResourceGeneratorService) Check(ctx context.Context, request *resource_service.HealthCheckRequest) (*resource_service.HealthCheckResponse, error) {
	// Create TODO: use proto validator instead of manual request validation https://github.com/bufbuild/protovalidate
	if len(request.GetService()) < 3 || len(request.GetService()) > 50 {
		return nil, fmt.Errorf("validation failure: service name should be between 3 and 50 characters")
	}
	deployment, err := s.Server.KubernetesClient.AppsV1().Deployments(namespace).Get(context.TODO(), request.Service, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("service error: %s", err)
	}
	// Print health information
	fmt.Printf("Deployment: %s\n", deployment.Name)
	fmt.Printf("Replicas: %d\n", deployment.Status.Replicas)
	fmt.Printf("Ready Replicas: %d\n", deployment.Status.ReadyReplicas)
	fmt.Printf("Available Replicas: %d\n", deployment.Status.AvailableReplicas)
	if deployment.Status.ReadyReplicas == deployment.Status.AvailableReplicas {
		return &resource_service.HealthCheckResponse{Status: resource_service.HealthCheckResponse_SERVING}, nil
	}
	return &resource_service.HealthCheckResponse{Status: resource_service.HealthCheckResponse_NOT_SERVING}, nil
}
