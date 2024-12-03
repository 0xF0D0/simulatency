package agent

import (
	"encoding/binary"
	"log"
	"strconv"

	podwatcher "github.com/0xf0d0/simulatency-agent/watcher/pod"
)

func (a * simulatencyAgent) updateEbpfMap(pod podwatcher.Pod, isAdd bool) error {
	if !isAdd {
		//Deletion
		err := a.ebpfMap.Delete(binary.NativeEndian.Uint32(pod.Ip))
		return err
	}

	//Addition

	if val, exist := pod.Annotations[annotationKey]; exist {
		key, err := strconv.Atoi(val)
		if err != nil {
			log.Println("string: " + val + "is not convertible to integer")
		}
		return a.ebpfMap.Put(uint32(key), binary.NativeEndian.Uint32(pod.Ip))
	}
	return nil
}