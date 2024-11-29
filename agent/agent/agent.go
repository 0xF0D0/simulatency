package agent

import (
	podwatcher "github.com/0xf0d0/simulatency-agent/watcher/pod"
)

var (
	_ Agent = &simulatencyAgent{}
)

type Agent interface {
	RunWatchLabel(string) error
}

type simulatencyAgent struct {
	k8sWatcher podwatcher.PodWatcher

}

func (a *simulatencyAgent) RunWatchLabel(label string) error {
	a.k8sWatcher.RunWatch(label) //TODO: may be add some ctx cancel logic?
	return nil
}

func (a * simulatencyAgent) updateEbpfMap(list []podwatcher.Pod) error {
	//Insert and Delete and Update ips
	
}