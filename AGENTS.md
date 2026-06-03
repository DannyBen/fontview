# Agent Notes

## Machine Constraints

This repository may live under `/vagrant`, where broad parallel file access can
hang the VM. Do not use parallel tool calls for commands that scan, test, build,
format, generate, render, use git, or otherwise walk many files under
`/vagrant`.

Run those commands one at a time.

## Project Shape

- `main.go` is the small binary entrypoint.
- `cmd/` owns CLI parsing and help text.
- `report/` owns font discovery, parsing, rendering, and serving.
- `report/templates/` contains the embedded HTML, CSS, and JavaScript assets.
- `op.conf` is the user-facing command catalog.
- `agent.op.conf` is for agent-specific repeatable checks only.

## Useful Commands

```bash
go test ./...
go run . --help
go run . --html
go run .
op check
```

`fontview.html` is generated output and should not be committed.
