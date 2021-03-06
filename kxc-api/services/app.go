package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/didil/kubexcloud/kxc-api/requests"
	"github.com/didil/kubexcloud/kxc-api/responses"

	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"
	"github.com/didil/kubexcloud/kxc-operator/controllers"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	validationutils "k8s.io/apimachinery/pkg/util/validation"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AppSvc interface {
	Create(ctx context.Context, projectName string, reqData *requests.CreateApp) error
	Update(ctx context.Context, projectName, appName string, reqData *requests.UpdateApp) error
	List(ctx context.Context, projectName string) (*responses.ListApp, error)
	Restart(ctx context.Context, projectName, appName string) error
}

type AppService struct {
	k8sSvc K8sSvc
}

func NewAppService(k8sSvc K8sSvc) *AppService {
	return &AppService{
		k8sSvc: k8sSvc,
	}
}

func (svc *AppService) validateCreateApp(reqData *requests.CreateApp) error {
	if reqData.Name == "" {
		return fmt.Errorf("name is required")
	}

	if errs := validationutils.IsDNS1123Label(reqData.Name); len(errs) > 0 {
		return fmt.Errorf("name: %s", strings.Join(errs, "."))
	}

	return nil
}

func (svc *AppService) Create(ctx context.Context, projectName string, reqData *requests.CreateApp) error {
	client := svc.k8sSvc.Client()

	err := svc.validateCreateApp(reqData)
	if err != nil {
		return fmt.Errorf("app invalid: %v", err)
	}

	app := &cloudv1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reqData.Name,
			Namespace: controllers.ProjectNamespaceName(projectName),
			Labels:    controllers.LabelsForApp(projectName, reqData.Name),
		},
		Spec: cloudv1alpha1.AppSpec{
			Replicas:   reqData.Replicas,
			Containers: []cloudv1alpha1.Container{},
		},
	}

	for _, c := range reqData.Containers {
		container := cloudv1alpha1.Container{
			Image:   c.Image,
			Name:    c.Name,
			Command: c.Command,
			Ports:   []cloudv1alpha1.Port{},
		}

		for _, p := range c.Ports {
			container.Ports = append(container.Ports, cloudv1alpha1.Port{
				Number:           p.Number,
				Protocol:         corev1.Protocol(p.Protocol),
				ExposeExternally: p.ExposeExternally,
			})
		}

		app.Spec.Containers = append(app.Spec.Containers, container)
	}

	err = client.Create(ctx, app)
	if err != nil {
		return fmt.Errorf("create app: %v", err)
	}
	return nil
}

func (svc *AppService) Update(ctx context.Context, projectName, appName string, reqData *requests.UpdateApp) error {
	client := svc.k8sSvc.Client()

	app := &cloudv1alpha1.App{}
	err := client.Get(ctx, types.NamespacedName{Name: appName, Namespace: controllers.ProjectNamespaceName(projectName)}, app)
	if err != nil {
		return fmt.Errorf("get app: %v", err)
	}

	app.Spec.Replicas = reqData.Replicas
	app.Spec.Containers = []cloudv1alpha1.Container{}

	for _, c := range reqData.Containers {
		container := cloudv1alpha1.Container{
			Image:   c.Image,
			Name:    c.Name,
			Command: c.Command,
			Ports:   []cloudv1alpha1.Port{},
		}

		for _, p := range c.Ports {
			container.Ports = append(container.Ports, cloudv1alpha1.Port{
				Number:           p.Number,
				Protocol:         corev1.Protocol(p.Protocol),
				ExposeExternally: p.ExposeExternally,
			})
		}

		app.Spec.Containers = append(app.Spec.Containers, container)
	}

	err = client.Update(ctx, app)
	if err != nil {
		return fmt.Errorf("update app: %v", err)
	}
	return nil
}

func (svc *AppService) List(ctx context.Context, projectName string) (*responses.ListApp, error) {
	cl := svc.k8sSvc.Client()

	appList := &cloudv1alpha1.AppList{}
	listOpts := []client.ListOption{
		client.InNamespace(controllers.ProjectNamespaceName(projectName)),
	}
	if err := cl.List(ctx, appList, listOpts...); err != nil {
		return nil, fmt.Errorf("failed to list apps: %v", err)
	}

	respData := &responses.ListApp{
		Apps: []responses.ListAppEntry{},
	}

	for _, app := range appList.Items {
		respData.Apps = append(respData.Apps, responses.ListAppEntry{
			Name:                app.Name,
			ExternalURL:         app.Status.ExternalURL,
			AvailableReplicas:   app.Status.AvailableReplicas,
			UnavailableReplicas: app.Status.UnavailableReplicas,
		})
	}

	return respData, nil
}

func (svc *AppService) Restart(ctx context.Context, projectName, appName string) error {
	client := svc.k8sSvc.Client()

	app := &cloudv1alpha1.App{}
	err := client.Get(ctx, types.NamespacedName{Name: appName, Namespace: controllers.ProjectNamespaceName(projectName)}, app)
	if err != nil {
		return fmt.Errorf("get app: %v", err)
	}

	if app.ObjectMeta.Annotations == nil {
		app.ObjectMeta.Annotations = map[string]string{}
	}

	app.ObjectMeta.Annotations[controllers.AppRestartAnnotationKey] = time.Now().UTC().Format(time.RFC3339)

	err = client.Update(ctx, app)
	if err != nil {
		return fmt.Errorf("restart app: %v", err)
	}
	return nil
}
