**Repository Overview**
- **Purpose:**: Simulation prototype implemented in Go under `code/backend/pkg`.
- **Primary package:**: `simulation` (all core types live in `code/backend/pkg`).
- **Current state:**: Small, partially-implemented package (several files contain skeletons).

**Big Picture / Architecture**
- **Core concepts:**: `Agent`, `Action`, `Environment`, `Simulation` (see `agent.go`, `action.go`, `environment.go`, `simulation.go`).
- **Data flow:**: Agents hold a `Position`/`Sprite` and choose `Action`s that are executed against an `Environment` (use `Execute(a *Agent, env *Environment)` signature in `action.go`).
- **Why this structure:**: The code separates entity data (structs) from decision behaviour (interfaces) to allow pluggable strategies and actions.

**Key files to read first**
- `code/backend/pkg/position.go`: position/vector utilities, `CreatePosition`, `CreateSprite`, `Position.DistanceTo`, `Sprite.IsColliding` (collision logic used by agents).
- `code/backend/pkg/action.go`: `Action` interface and parameter struct; includes `Execute` and `evaluateUtility` (utility-based action selection pattern).
- `code/backend/pkg/agent.go`: agent interface and `Agent` struct skeleton — this is intentionally minimal and currently incomplete.
- `code/backend/pkg/environment.go`: environment container (width/height and agent list).
- `code/backend/pkg/simulation.go`: simulation loop state holder (maxSteps, currentStep, agents, environment).

**Project-specific conventions and patterns**
- **Unexported core types:** Many domain types are lowercase (e.g., `human`), so prefer keeping internal package scope unless you intentionally create a public API.
- **Constructor-style helpers:** Factory functions are named `CreateX` (e.g., `CreatePosition`, `CreateSprite`, `CreateVector`) — follow this convention for new constructors.
- **Method receivers:** Use pointer receivers when mutating (e.g., `(*Position).Move`, `(*Position).SetPosition`) and value receivers when returning immutable data (`(Position) GetPosition`).
- **Utility-based actions:** `Action` exposes an `evaluateUtility` method that returns a `float64` — the codebase favors utility scoring to pick actions.

**Integration points & dependencies**
- **Go module:** `code/go.mod` sets `module projet` and `go 1.25.1`. Don't change the module path lightly.
- **No external services:** There are no HTTP servers, databases, or external APIs in the repository; most work is internal library code.

**Developer workflows (commands)**
- **Build everything:** `cd code && go build ./...`
- **Run vet & tests:** `cd code && go vet ./...` and `cd code && go test ./...` (currently there may be no tests).
- **Format:** `gofmt -w .` (run from `code` to reformat Go files).
- **Linting (optional):** `golangci-lint run` if the project maintainer has it installed.

**How an AI agent should edit this repo**
- **Keep package scope consistent:** Changes to types that alter visibility (lowercase → exported) must be intentional and coordinated; prefer adding exported wrappers rather than changing many existing names.
- **Follow existing constructors and naming:** Use `CreateX` constructors and existing method names (`DistanceTo`, `IsColliding`).
- **Preserve small, focused changes:** Files are small and incomplete — make minimal, well-scoped edits and run `go build ./...` to verify compilation.
- **When adding new behaviour:** Implement new `Action` types by creating a struct that satisfies `Action` (implement `Execute` and `evaluateUtility`) and add tests under the `code` directory.

**Examples (copy-paste friendly)**
- Add a new action skeleton:

```go
type ForageAction struct{ ParametesAction }

func (f *ForageAction) Execute(a *Agent, env *Environment) {
    // implement state changes
}

func (f *ForageAction) evaluateUtility(a *Agent, env *Environment) float64 {
    return 0.0
}
```

- Use `Position` helper:

```go
pos := CreatePosition(1.0, 2.0)
v := CreateVector(0.5, -0.2)
pos.Move(v)
```

**What to watch for / current TODOs**
- `agent.go`, `action.go`, and `ressources.go` contain incomplete implementations — expect build errors until these are implemented.
- The codebase encodes numeric behaviour with `float64` for positions and energy/hunger values — keep that type consistent.

**If unclear / missing information**
- Ask for the intended public API for `Agent` and `Action` selection logic (e.g., how actions are prioritized in simulation loop).
- Confirm whether tests should use `package simulation` or `simulation_test` to access internal types.

Please review this file for anything you want emphasized or adjusted (e.g., add CI commands, test conventions, or preferred linters).  