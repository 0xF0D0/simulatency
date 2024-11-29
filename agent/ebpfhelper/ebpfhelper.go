package ebpfhelper

import (
	"log"
	"github.com/cilium/ebpf/rlimit"	
)

func init() {
	if err := rlimit.RemoveMemlock(); err != nil { 
		log.Fatal("Removing memlock:", err)
	}
}


func Attach() {

}