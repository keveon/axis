// Package sdk provides the shared client library for Axis components.
//
// Independent binaries (Listener, Collector, Parser, Sink, Relay) use the SDK to:
//   - Register with the Axis Server and maintain heartbeats
//   - Push raw data to the Server (Sources) or subscribe to events (Sinks/Parsers)
//   - Manage checkpoints with CAS semantics
//   - Watch configuration changes
//
// The SDK also provides framework-level abstractions (RunCollector, RunParser,
// RunSink, etc.) that handle the full component lifecycle so developers only
// implement the business logic interface.
package sdk
