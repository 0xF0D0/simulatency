package k8s

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sClient interface {
	CoordinationV1() coordinationv1.CoordinationV1Interface

	CreateStatefulset(ctx context.Context, statefulset *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	DeleteStatefulset(ctx context.Context, name string) error
	GetStatefulset(ctx context.Context, name string) (*appsv1.StatefulSet, error)
	GetDeployment(ctx context.Context, name string) (*appsv1.Deployment, error)
	ApplyScaleStatefulset(ctx context.Context, name string, replicas int32) (*appsv1.StatefulSet, error)
	CreateService(ctx context.Context, service *corev1.Service) (*corev1.Service, error)
	DeleteService(ctx context.Context, name string) error
	DeletePersistentVolumeClaim(ctx context.Context, name string) error
	DeletePod(ctx context.Context, name string) error
	ListPods(ctx context.Context) (*corev1.PodList, error)
	CreateJob(ctx context.Context, job *batchv1.Job) (*batchv1.Job, error)
	DeleteJob(ctx context.Context, name string) error
	CreatePvc(ctx context.Context, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error)
	DeletePvc(ctx context.Context, name string) error
	GetImage(ctx context.Context, namespace, podName string) (string, error)
	GetImageFromDeployment(ctx context.Context, namespace, podName string) (string, error)
}

type k8sClient struct {
	client    *kubernetes.Clientset
	namespace string
}

func NewClient(namespace string) (K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	config.QPS = 100
	config.Burst = 200
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	clientset.CoordinationV1()

	return &k8sClient{
		client:    clientset,
		namespace: namespace,
	}, nil
}

func (c *k8sClient) CoordinationV1() coordinationv1.CoordinationV1Interface {
	return c.client.CoordinationV1()
}
func (c *k8sClient) CreateStatefulset(ctx context.Context, statefulset *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return c.client.AppsV1().StatefulSets(c.namespace).Create(ctx, statefulset, metav1.CreateOptions{})
}

func (c *k8sClient) DeleteStatefulset(ctx context.Context, name string) error {
	return c.client.AppsV1().StatefulSets(c.namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *k8sClient) GetStatefulset(ctx context.Context, name string) (*appsv1.StatefulSet, error) {
	return c.client.AppsV1().StatefulSets(c.namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *k8sClient) GetDeployment(ctx context.Context, name string) (*appsv1.Deployment, error) {
	return c.client.AppsV1().Deployments(c.namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *k8sClient) ApplyScaleStatefulset(ctx context.Context, name string, replicas int32) (*appsv1.StatefulSet, error) {
	sts, err := c.GetStatefulset(ctx, name)
	if err != nil {
		return nil, err
	}

	*sts.Spec.Replicas = replicas
	return c.client.AppsV1().StatefulSets(c.namespace).Update(ctx, sts, metav1.UpdateOptions{})
}

func (c *k8sClient) ListPods(ctx context.Context) (*corev1.PodList, error) {
	return c.client.CoreV1().Pods(c.namespace).List(ctx, metav1.ListOptions{})
}

func (c *k8sClient) CreateService(ctx context.Context, service *corev1.Service) (*corev1.Service, error) {
	return c.client.CoreV1().Services(c.namespace).Create(ctx, service, metav1.CreateOptions{})
}

func (c *k8sClient) DeleteService(ctx context.Context, name string) error {
	return c.client.CoreV1().Services(c.namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *k8sClient) DeletePersistentVolumeClaim(ctx context.Context, name string) error {
	return c.client.CoreV1().PersistentVolumeClaims(c.namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *k8sClient) DeletePod(ctx context.Context, name string) error {
	return c.client.CoreV1().Pods(c.namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *k8sClient) CreateJob(ctx context.Context, job *batchv1.Job) (*batchv1.Job, error) {
	return c.client.BatchV1().Jobs(c.namespace).Create(ctx, job, metav1.CreateOptions{})
}

func (c *k8sClient) DeleteJob(ctx context.Context, name string) error {
	return c.client.BatchV1().Jobs(c.namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *k8sClient) CreatePvc(ctx context.Context, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	return c.client.CoreV1().PersistentVolumeClaims(c.namespace).Create(ctx, pvc, metav1.CreateOptions{})
}

func (c *k8sClient) DeletePvc(ctx context.Context, name string) error {
	return c.client.CoreV1().PersistentVolumeClaims(c.namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *k8sClient) GetImage(ctx context.Context, namespace, podName string) (string, error) {
	pod, err := c.client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return pod.Spec.Containers[0].Image, nil
}

func (c *k8sClient) GetImageFromDeployment(ctx context.Context, namespace, podName string) (string, error) {
	deployment, err := c.client.AppsV1().Deployments(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return deployment.Spec.Template.Spec.Containers[0].Image, nil
}
