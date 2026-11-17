# Validation record

Executed in the authoring environment:

- `go test ./...`: passed.
- `kubectl kustomize`: passed for Module 01, Module 02 base, dev and prod.
- Local Markdown links: checked programmatically; all course links resolve.
- Numbered lab questions: checked against instructor answers; all covered.
- Workflow location: executable YAML is at the golang-app repository root `.forgejo/workflows/course-ci.yaml`.
- GitOps paths: both point to existing overlays; repository URL intentionally requires student substitution.
- Ownership: Argo targets course-dev, Flux course-prod; handoff removes Argo Application without cascading deletion.
- Node configuration: one control-plane and two workers; kind name course generates the three documented names.

Not validated end to end: Docker daemon access is denied by the authoring sandbox. Image build and kind discovery were attempted and failed at socket access. No course cluster was created. The existing context was raspberry-pi and was not modified. GitHub publication, remote image pulls, Argo/Flux installation, rolling updates and live exercises require execution by the instructor/student. PyYAML is unavailable; Kustomize parsed application YAML, but a separate complete YAML syntax/schema validation was not executed. Version pins are historical teaching baselines, not current production recommendations. Validate upgrades and controller compatibility before classroom delivery.

## Simplified classroom demos

Kustomize has a 30-minute teaching allocation with a short three-question lab. Longer labs and solutions remain as extended-lab.md and extended-lab-solutions.md. Kustomize manifests and GitOps paths are unchanged; all three Kustomize configurations rendered successfully again. All local course links and short-demo answer coverage passed checks.

## Three-tier capstone replacement

The final project now requires frontend + Go API + persistent database, a Helm or Kustomize package, one GitOps controller, and a ConfigMap-driven visual change. The complete instructor implementation, tests and walkthrough are under 05-capstone/instructor-solutions/. Go race tests, vet, build, Kustomize rendering, plain YAML parsing and shell syntax checks passed. Docker/live database and cluster tests, Helm rendering and registry publication remain blocked or unexecuted; see the instructor solution's VALIDATION.md for the exact scope.

## Instructor solution layout refactor

The capstone solution now has independent api/ and frontend/ modules, deployment/{kustomize,helm,gitops}, and docs/ guides. Both services passed race tests, vet and Linux arm64 compilation after moving. Kustomize output is unchanged. GitOps paths, root CI build contexts, local documentation links and scripts were updated and checked. No live cluster or persisted data was modified. See 05-capstone/instructor-solutions/docs/VALIDATION.md for integration limitations.
