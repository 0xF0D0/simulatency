package pod

import (
	"net"
)

type Pod struct {
	Name string
	Ip net.IP
	Labels map[string] string
	Annotations map[string] string
}

type PodWatcher interface {
	RunWatch(label string)
	SetHandler(func ([]Pod) error)
}

