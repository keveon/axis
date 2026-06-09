// Package controllermanager implements the Axis Controller Manager.
//
// It hosts all built-in controllers — PipelineController, StoreController,
// and ComponentController — each running an independent watch-reconcile loop
// against etcd to drive actual state toward desired state.
package controllermanager
