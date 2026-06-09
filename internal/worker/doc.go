// Package worker implements the generic Axis Worker logic.
//
// A Worker is a single binary that serves all Store instances. On startup,
// it reads its assigned Store CRD from the Axis Server, connects to the
// target database, applies schema migrations, and exposes gRPC read/write
// interfaces for Sinks and Collectors.
//
// Multiple Store instances (PostgreSQL, MySQL, TimescaleDB) can coexist on
// different machines, each backed by the same axis binary with the worker role
// and different CRD configurations.
package worker
