package pod

import (
	"net"
)

type Pod struct {
	name string
	ip net.IP
	labels map[string] string
	annotations map[string] string
}

type PodWatcher interface {
	RunWatch(label string)
	SetHandler(func ([]Pod) error)
}

