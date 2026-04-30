# care-giver-api

AWS Lambda function behind API Gateway that serves as the core REST API for CareToSher. Deployed via AWS SAM. All routes require a Cognito JWT except `POST /user`.

## Commands

```bash
make test           # unit tests with coverage
make test-report    # open coverage HTML in browser
make lint           # golangci-lint
make invoke         # build + run a single Lambda invocation locally (default: events/event.json)
make start-api      # build + run API Gateway locally (requires Docker + care-giver-infra_default network)
```

Override the local event file: `make invoke EVENT=events/some_other_event.json`

## Deployment

**Never run `sam deploy` manually.** CI handles all deployments:

- **Every PR opened/updated** → auto-deploys to `dev` (`api-dev.caretosher.com`)
- **Merge to `main`** → auto-deploys to `prod` (`api.caretosher.com`)

Both pipelines use OIDC (no stored AWS keys). `confirm_changeset` is `false` in CI; it is `true` in `samconfig.toml` for local use as a safety net.

## Adding a New Endpoint

Four things must happen together — missing any one will break the route:

1. **Handler function** — new file in `internal/handlers/` following the existing pattern (see below).
2. **Registry entry** — add `{"/your/path", http.MethodXxx}: HandleYourFunc` to `handlersMap` in `internal/handlers/handlers.go`.
3. **SAM route** — add the corresponding `Events` entry and any IAM permissions to `template.yaml`.
4. **Local test event** — add a JSON file under `events/` so the route can be invoked with `make invoke`.

## Handler Pattern

Every handler has the same signature:

```go
func HandleFoo(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error)
```

All dependencies (repos, config, logger) come through `HandlerParams` — handlers have no package-level state. Follow this sequence inside each handler:

1. Parse + validate the request body with `readRequestBody` (or `validatePathParameters` / `validateQueryParameters` for non-body inputs).
2. Look up the requesting user from `params.UserRepo`.
3. **Check authorization** with `relationship.IsACareGiver()` before touching any receiver data. Every handler that acts on a receiver must do this — skip it and you have a broken access control bug.
4. Perform the business logic.
5. Return via `response.FormatResponse(payload, http.StatusOK)` or one of the error constructors (`CreateBadRequestResponse`, `CreateInternalServerErrorResponse`, `CreateAccessDeniedResponse`).

## Key Conventions

- `readRequestBody` uses `json.Decoder` with `DisallowUnknownFields()` — request structs must be an exact match for the JSON body. Unknown fields return 400.
- Timestamps must be RFC3339; use `validateTimestamps` which also checks end ≥ start.
- ID path parameters are validated against the format `Prefix#uuid` (e.g. `Receiver#...`). The `validatePathParameters` helper handles URL-decoded `%23` → `#` automatically.
- Log keys are constants in `care-giver-golang-common/pkg/log` — never use raw strings as zap field keys.
- `removePathPrefix` strips `/Stage` and `/Prod` from resource paths. If a new API Gateway stage is added, update that function.

## Testing

- Table-driven tests with `testify/assert`, one `*_test.go` per handler file.
- Repositories are injected via their `*Provider` interfaces — use `dynamo.Mock` from the common library to stub DynamoDB. Never hit a real table in unit tests.
- `mock_test.go` in `internal/handlers/` contains shared mock implementations reused across handler tests.

## Dependency on care-giver-golang-common

The common library is vendored (`vendor/`). After bumping the version in `go.mod`, run:

```bash
go get github.com/care-giver-app/care-giver-golang-common@vX.Y.Z
go mod vendor
```

Commit both `go.mod`, `go.sum`, and the updated `vendor/` directory together.
