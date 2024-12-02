package pod

import (
	"context"
	"fmt"
	"net"

	"github.com/0xf0d0/simulatency-agent/client/k8s"
	"k8s.io/apimachinery/pkg/watch"
	corev1 "k8s.io/api/core/v1"
)

var (
	_ PodWatcher = (*podWatcher)(nil)
)

type Pod struct {
	Name        string
	Ip          net.IP
	Labels      map[string]string
	Annotations map[string]string
}

type PodWatcher interface {
	RunWatch(label string)
	SetHandler(func(p Pod, isAdd bool) error) //Do net set after watch
}

type podWatcher struct {
	k8sClient k8s.K8sClient
	handlers []func(Pod, bool) error
}

// RunWatch implements PodWatcher.
func (p *podWatcher) RunWatch(label string) {
	watcher, err := p.k8sClient.ListPod(context.Background(), label)
	if err != nil {
		panic(err)
	}
	// Acts as a main thread
	for e := range watcher.ResultChan() {
		switch e.Type {
		case watch.Added, watch.Deleted:
			item := e.Object.(*corev1.Pod)
			for _, h := range p.handlers {
				ip := net.ParseIP(item.Status.PodIP)
				if ip == nil {
					fmt.Println("pod ip parse failed: "+ item.Status.PodIP)
					continue
				}
				go h(Pod{
					Name:        item.Name,
					Ip:          ip,
					Labels:      item.Labels,
					Annotations: item.Annotations,
				}, e.Type == watch.Added)
			}
		case watch.Bookmark, watch.Error, watch.Modified:
			continue
		default:
			panic(fmt.Sprintf("unexpected watch.EventType: %#v", e.Type))
	}
	}
}

// SetHandler implements PodWatcher.
func (p *podWatcher) SetHandler(h func(Pod, bool) error) {
	p.handlers = append(p.handlers, h)
}
