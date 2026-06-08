## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

## 5. Repository Architecture & Boundaries

The project follows **Hexagonal Architecture (Ports and Adapters / Clean Architecture)**. Dependencies flow strictly **inwards** toward the core domain.

### Directory Structure & Responsibilities:
- `cmd/`: Application entry point (`main.go`). Initializes configurations, connection pools, and runs database migrations.
- `internal/applications/`: Manages runtime apps (e.g. HTTP server, background scraper). Wires dependencies and manages runtime lifecycles.
- `internal/command/`: Domain orchestrators/use cases containing application business workflows. Decoupled from transport (HTTP) and DB drivers.
- `internal/domain/`: Core models and ports/interfaces. **Zero** imports from adapters, commands, or transport layers.
- `internal/adapter/pgx/`: Postgres database adapters implementing domain repository ports. SQL files must be colocated in a `queries/` subdirectory inside each DB package, loaded via `//go:embed`, and run using the separated `readDB` and `writeDB` `pgxpool.Pool` connections.
- `internal/adapter/collector/`: Outbound provider integrations (e.g. status page scrapers) implementing `StatusPageProvider`.
- `internal/port/`: Inbound controllers (e.g. OpenAPI strict HTTP handlers) that route external requests to commands.
- `internal/common/`: Common helpers (status normalizer, Snowflake IDs, application error formats, local queue wrappers).

### Rules and Boundaries to Maintain:
- **No Outward Domain Imports**: Core domain code must never depend on database packages (`internal/adapter/pgx/...`), commands, or handlers.
- **Constructor Injection**: All repositories/providers must be injected as interfaces into commands and applications.
- **DB & Transport Abstraction**: No raw database client structs (`pgx`) or HTTP framework imports are allowed within command handlers.
- **Colocated DB SQL**: Always colocate SQL files in `queries/` under the corresponding DB adapter package and embed them in Go code.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.