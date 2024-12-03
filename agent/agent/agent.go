package agent

import (
	"github.com/0xf0d0/simulatency-agent/ebpf"
	"github.com/0xf0d0/simulatency-agent/util"
	podwatcher "github.com/0xf0d0/simulatency-agent/watcher/pod"
)

const (
	annotationKey = "simulatency/longitude"
)

var (
	_ Agent = &simulatencyAgent{}
)

type Agent interface {
	RunWatchLabel(string) error
}

type simulatencyAgent struct {
	k8sWatcher podwatcher.PodWatcher
	ebpfMap ebpf.EbpfMap
	distanceMap util.CombinationMap
}

func InitializeAgent(watcher podwatcher.PodWatcher, ebpfMap ebpf.EbpfMap) (Agent) {
	a := &simulatencyAgent{
		k8sWatcher:  watcher,
		ebpfMap:     ebpfMap,
		distanceMap: util.NewCombinationMap(),
	}
	watcher.SetHandler(a.updateEbpfMap)
	watcher.SetHandler(a.injectDelays)
	return a
}

func (a *simulatencyAgent) RunWatchLabel(label string) error {
	a.k8sWatcher.RunWatch(label) //TODO: may be add some ctx cancel logic?
	return nil
}