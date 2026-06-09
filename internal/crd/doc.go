// Package crd defines all Axis resource types (CRDs) as Go structs.
//
// Each resource follows the spec/status pattern inspired by Kubernetes:
//   - spec: desired state declared by the user
//   - status: observed state maintained by controllers
//
// Resource types include Protocol, Pipeline, Store, Listener, Collector,
// Parser, Sink, and Relay.
package crd
