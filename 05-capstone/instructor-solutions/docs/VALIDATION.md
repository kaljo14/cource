# Instructor solution validation

## Refactor verification

After splitting the services into independent api/ and frontend/ modules:

- `go test -race ./...` and `go vet ./...` passed separately in both directories.
- Both services cross-compiled as static Linux arm64 binaries with their Dockerfile linker flags.
- The Kustomize demo render was byte-for-byte identical before and after moving it into deployment/.
- YAML parsing, all course local Markdown links, GitOps source paths, expanded CI matrix paths, executable permissions and shell syntax passed checks.
- The existing yellow theme and saved demo.yaml were preserved. The snapshot is now in deployment/rendered/ and is not a GitOps source.
- No live cluster resources were changed by this refactor. Docker/Helm and live integration remain subject to the limitations below.

The service directories are separate Go modules. Run Go checks inside each one, not from instructor-solutions/.

## Earlier checks during implementation

- `go test -race ./...`: passed for API request validation, note creation/listing using a fake store, database-unavailable responses, frontend proxy routing, theme file rereads, static page delivery and invalid-color fallback.
- `go vet ./...` and `go build ./...`: passed.
- Both components cross-compiled successfully as static Linux arm64 binaries with CGO disabled and the Dockerfile linker flags. This checks compilation, not execution inside the container.
- `go mod tidy` resolved the pinned lib/pq dependency from the local cache and generated go.sum. No database driver source was invented or vendored manually.
- `kubectl kustomize deployment/kustomize/overlays/demo`: passed; rendered all resources into capstone.
- Ruby YAML parsing: passed for plain Kubernetes manifests, GitOps manifests, chart metadata/values, rendered Kustomize output and root Actions workflows. This checks syntax, not Kubernetes API admission or Helm template evaluation.
- Shell syntax: all three scripts passed `bash -n`.
- Local documentation links and student Q1–Q6 reference-answer coverage checked.

The initial frontend test attempted a local HTTP test server, which this sandbox disallowed. It was replaced with an injected HTTP transport that tests reverse-proxy requests and failures without opening a listening socket. These remain unit tests, not live network integration tests.

## Not executed here

Docker build was attempted and failed because the sandbox cannot access the Docker daemon socket. No capstone cluster or PostgreSQL container was started. Database initialization, actual SQL persistence, local-path volume permissions, image pulls, live color propagation, rolling updates and GitOps adoption still require the local instructor run.

Helm is not installed in this environment, so the chart has not been rendered, linted, packaged or installed by Helm. Argo/Flux CRDs were not installed or server-validated. GitHub Actions and GHCR publication were not executed against an account. No existing cluster was modified.

Go checks ran with the installed host compiler; the pinned Go 1.23.6 container build has not been executed. Image/tool pins are historical teaching baselines, not a security-support guarantee. Review version updates before production use.

## Local acceptance checklist

1. Follow START.md on kind-capstone; both app rollouts and db-0 must become ready.
2. Run scripts/smoke.sh with the frontend forward active. Save/list a note.
3. Delete db-0, wait for recovery and verify that note remains.
4. Change the chosen package's background color; verify /config and the browser change without an image rebuild.
5. Follow one GITOPS.md route; verify source revision, health and replica drift correction.
6. Promote API v2, inject a missing tag, then revert through Git and verify recovery.
7. For Helm, execute lint/template/package and perform its live deployment on a fresh prepared cluster rather than overwriting a Kustomize installation.
8. Export any needed data and use the documented cleanup. Confirm the course cluster remains untouched.
