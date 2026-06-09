// Package agent implements the Axis Agent.
//
// The Agent runs on each machine, registers with the Axis Server, and manages
// the lifecycle of systemd services for deployed components (Listener, Parser,
// Sink, Worker, Relay). It receives placement decisions from the Scheduler
// and reports machine status back to the Server.
package agent
