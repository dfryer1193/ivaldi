# Mjolnir Roadmap / Upstream Tasks

This document contains a checklist of features and improvements identified during the development of the `ivaldi` project generator that should be considered for upstreaming into the `mjolnir` library.

## Infrastructure & Routing
- [ ] **Built-in Standard Health Checks:**
  - Add a built-in standard `/health` endpoint to `mjolnirRouter.go` that returns a simple `200 OK` JSON response.
  - Add a built-in standard `/live` endpoint for container liveness probes.
- [ ] **Readiness Check API:**
  - Design and implement a readiness registration API within `mjolnir`. This would allow applications to register dependency checks (e.g., `mjolnir.RegisterReadyCheck("database", db.Ping)`) rather than manually wiring up the `/ready` route in every `main.go`.

## Logging
- [ ] **Extract `logx` Package:**
  - Currently, `mjolnir` sets up a global `zerolog` configuration directly in the router initialization (`mjolnirRouter.go`).
  - Extract this `zerolog` configuration into a dedicated `logx` package within `mjolnir/utils`.
  - This allows non-HTTP binaries (like background workers) to easily import and use the exact same logging format without having to duplicate the `zerolog.ConsoleWriter` setup code.
