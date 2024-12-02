package agent

import (
	"encoding/binary"
	"log"
	"strconv"

	"github.com/0xf0d0/simulatency-agent/ebpf"
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
}

func InitializeAgent(watcher podwatcher.PodWatcher, ebpfMap ebpf.EbpfMap) (Agent) {
	a := &simulatencyAgent{
		k8sWatcher: watcher,
		ebpfMap:    ebpfMap,
	}
	watcher.SetHandler(a.updateEbpfMap)
	return a
}

func (a *simulatencyAgent) RunWatchLabel(label string) error {
	a.k8sWatcher.RunWatch(label) //TODO: may be add some ctx cancel logic?
	return nil
}

func (a * simulatencyAgent) updateEbpfMap(list []podwatcher.Pod) error {
	//Insert and Delete and Update ips
	for _, pod := range list  {
		if val, exist := pod.Annotations[annotationKey]; exist {
			key, err := strconv.Atoi(val)
			if err != nil {
				log.Println("string: " + val + "is not convertible to integer")
				continue
			}
			ipValue := binary.NativeEndian.Uint32(pod.Ip)
			a.ebpfMap.Put(uint32(key), ipValue)
		}
	}

	return nil
}