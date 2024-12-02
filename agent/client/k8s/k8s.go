package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var _ K8sClient = (*k8sClient)(nil)

type K8sClient interface {
	ListPod(ctx context.Context, label string) (watch.Interface, error)
}

type k8sClient struct {
	client *kubernetes.Clientset
}

func NewClient() (K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	config.QPS = 100
	config.Burst = 200

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &k8sClient{
		client: clientset,
	}, nil
}

// listPod implements K8sClient.
func (k *k8sClient) ListPod(ctx context.Context, label string) (watch.Interface, error) {
	t := true
	return k.client.CoreV1().Pods("").Watch(ctx, metav1.ListOptions{
		LabelSelector:        label,
		Watch:                true,
		SendInitialEvents:    &t,
	})
}