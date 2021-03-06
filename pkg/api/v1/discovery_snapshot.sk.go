// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"fmt"

	"github.com/solo-io/solo-kit/pkg/utils/hashutils"
	"go.uber.org/zap"
)

type DiscoverySnapshot struct {
	Pods   PodsByNamespace
	Meshes MeshesByNamespace
}

func (s DiscoverySnapshot) Clone() DiscoverySnapshot {
	return DiscoverySnapshot{
		Pods:   s.Pods.Clone(),
		Meshes: s.Meshes.Clone(),
	}
}

func (s DiscoverySnapshot) Hash() uint64 {
	return hashutils.HashAll(
		s.hashPods(),
		s.hashMeshes(),
	)
}

func (s DiscoverySnapshot) hashPods() uint64 {
	return hashutils.HashAll(s.Pods.List().AsInterfaces()...)
}

func (s DiscoverySnapshot) hashMeshes() uint64 {
	return hashutils.HashAll(s.Meshes.List().AsInterfaces()...)
}

func (s DiscoverySnapshot) HashFields() []zap.Field {
	var fields []zap.Field
	fields = append(fields, zap.Uint64("pods", s.hashPods()))
	fields = append(fields, zap.Uint64("meshes", s.hashMeshes()))

	return append(fields, zap.Uint64("snapshotHash", s.Hash()))
}

type DiscoverySnapshotStringer struct {
	Version uint64
	Pods    []string
	Meshes  []string
}

func (ss DiscoverySnapshotStringer) String() string {
	s := fmt.Sprintf("DiscoverySnapshot %v\n", ss.Version)

	s += fmt.Sprintf("  Pods %v\n", len(ss.Pods))
	for _, name := range ss.Pods {
		s += fmt.Sprintf("    %v\n", name)
	}

	s += fmt.Sprintf("  Meshes %v\n", len(ss.Meshes))
	for _, name := range ss.Meshes {
		s += fmt.Sprintf("    %v\n", name)
	}

	return s
}

func (s DiscoverySnapshot) Stringer() DiscoverySnapshotStringer {
	return DiscoverySnapshotStringer{
		Version: s.Hash(),
		Pods:    s.Pods.List().NamespacesDotNames(),
		Meshes:  s.Meshes.List().NamespacesDotNames(),
	}
}
