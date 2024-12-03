package agent

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"

	podwatcher "github.com/0xf0d0/simulatency-agent/watcher/pod"
)

//TODO: add modified case
func (a * simulatencyAgent) injectDelays(pod podwatcher.Pod, isAdd bool) error {
	val, exist := pod.Annotations[annotationKey]
	if !exist {
		return fmt.Errorf("longitude not defined")
	}

	longitude, err := strconv.Atoi(val)
	if err != nil {
		log.Println("string: " + val + "is not convertible to integer")
	}


	if !isAdd {
		removedDistances := a.distanceMap.Remove(longitude)
		// Do something with delete distances!
		return err
	} else {
		addedDistances := a.distanceMap.Insert(longitude)
		// Do something with added distances!
	}
	return nil
}
