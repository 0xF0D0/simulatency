package ebpf

import (
	"log"
	"net"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

var (
	_ EbpfMap = &SimulatencyEBPFModule{}
)

func init() {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("Removing memlock:", err)
	}
}

type EbpfMap interface {
	Lookup(uint32) (uint32, error)
	Put(k, v uint32) error
}

type SimulatencyEBPFModule struct {
	objs simulatencyObjects
	link link.Link
}

func (m *SimulatencyEBPFModule) Attach(netInterfaceName string) {

	var objs simulatencyObjects
	if err := loadSimulatencyObjects(&objs, nil); err != nil {
		log.Fatal("Loading eBPF objects:", err)
	}
	m.objs = objs

	iface, err := net.InterfaceByName(netInterfaceName)
	if err != nil {
		log.Fatalf("Getting interface %s: %s", netInterfaceName, err)
	}

	link, err := link.AttachTCX(link.TCXOptions{
		Interface: iface.Index,
		Program:   objs.ClassifyTcp,
		Attach:    ebpf.AttachTCXEgress,
	})

	if err != nil {
		log.Fatalf("could not attach TCx program: %s", err)
	}
	m.link = link
}

func (m *SimulatencyEBPFModule) Close() {
	err := m.objs.Close()
	log.Println(err)
	err = m.link.Close()
	log.Println(err)
}

// Lookup implements EbpfMap.
func (m *SimulatencyEBPFModule) Lookup(key uint32) (uint32, error) {
	var val uint32
	err := m.objs.IpTagMap.Lookup(key, &val)
	return val, err
}

// Put implements EbpfMap.
func (m *SimulatencyEBPFModule) Put(k uint32, v uint32) error {
	return m.objs.IpTagMap.Put(k, v)
}