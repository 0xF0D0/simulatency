// Code generated by bpf2go; DO NOT EDIT.
//go:build mips || mips64 || ppc64 || s390x

package ebpf

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

// loadSimulatency returns the embedded CollectionSpec for simulatency.
func loadSimulatency() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_SimulatencyBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load simulatency: %w", err)
	}

	return spec, err
}

// loadSimulatencyObjects loads simulatency and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*simulatencyObjects
//	*simulatencyPrograms
//	*simulatencyMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadSimulatencyObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadSimulatency()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// simulatencySpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type simulatencySpecs struct {
	simulatencyProgramSpecs
	simulatencyMapSpecs
}

// simulatencySpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type simulatencyProgramSpecs struct {
	ClassifyTcp *ebpf.ProgramSpec `ebpf:"classify_tcp"`
}

// simulatencyMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type simulatencyMapSpecs struct {
	IpTagMap *ebpf.MapSpec `ebpf:"ip_tag_map"`
}

// simulatencyObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadSimulatencyObjects or ebpf.CollectionSpec.LoadAndAssign.
type simulatencyObjects struct {
	simulatencyPrograms
	simulatencyMaps
}

func (o *simulatencyObjects) Close() error {
	return _SimulatencyClose(
		&o.simulatencyPrograms,
		&o.simulatencyMaps,
	)
}

// simulatencyMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadSimulatencyObjects or ebpf.CollectionSpec.LoadAndAssign.
type simulatencyMaps struct {
	IpTagMap *ebpf.Map `ebpf:"ip_tag_map"`
}

func (m *simulatencyMaps) Close() error {
	return _SimulatencyClose(
		m.IpTagMap,
	)
}

// simulatencyPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadSimulatencyObjects or ebpf.CollectionSpec.LoadAndAssign.
type simulatencyPrograms struct {
	ClassifyTcp *ebpf.Program `ebpf:"classify_tcp"`
}

func (p *simulatencyPrograms) Close() error {
	return _SimulatencyClose(
		p.ClassifyTcp,
	)
}

func _SimulatencyClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed simulatency_bpfeb.o
var _SimulatencyBytes []byte
