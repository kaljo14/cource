# Student requirements

## 1. Build a small, complete application

Implement a noticeboard or an equivalent simple application with the same observable behavior:

- A browser frontend displays records and submits a new record through an API.
- A Go API implements GET /api/notes and POST /api/notes (or documented equivalent routes), validates input and uses parameterized database operations.
- A separate database stores records; use PostgreSQL for the simplest match to the reference solution. Data must not live only in a process variable or frontend file.
- Include API liveness and readiness endpoints. Readiness must report database unavailability; liveness should not restart an otherwise healthy API process merely because the database is temporarily down.
- Include tests for valid input, invalid input and a dependency failure. Keep UI scope small: one page, one form and a list are sufficient. Authentication and complex styling are not required; expose the demo only on localhost.

## 2. Build and publish

- Provide separate API and frontend container images and reproducible Dockerfiles. Use a multi-stage build for the Go API and run application containers as non-root.
- Use the official database image; do not write your own database server.
- STACKIT Git Actions must test/build and publish API/frontend images to GHCR on trusted branch pushes. Pull requests must not publish. Use traceable tags and minimal permissions.
- Document actual image references. Never commit tokens or passwords.

## 3. Create your own kind cluster

Create a separate cluster named capstone using the Module 00 three-node configuration. Expected node names are capstone-control-plane, capstone-worker and capstone-worker2. Keep the earlier course cluster separate, or delete it deliberately after saving any work. Two concurrent clusters require additional host resources.

Use namespace capstone for the application. Provide Services for the frontend, API and database. Access the frontend through localhost port-forwarding; do not expose the database or GitOps administrative UI publicly.

Run at least two API and two frontend replicas; one database replica is sufficient. Include probes, requests/limits and appropriate application security contexts. Persist database data using a PVC and demonstrate that a database Pod replacement retains a record. kind node-local storage is a teaching facility, not a highly available database or backup.

## 4. Package all three tiers: choose one

**Kustomize:** a reusable base and at least one environment overlay containing the application resources, Services, PVC and non-secret configuration. Render it before applying. Avoid copied environment manifests.

**Helm:** one chart containing the three tiers, Services, PVC and ConfigMap, with values for image references, application replica counts and background color. Lint, render and package the chart. State chart version separately from image versions. Publishing the chart to GHCR is optional; publishing the application images is required.

Create the database Secret separately and document how to supply it. Never package a real credential into a chart or commit a plaintext Secret. Choose one packaging route rather than combining both unnecessarily.

## 5. GitOps: choose one controller

Use Argo CD or Flux on the capstone cluster to reconcile your chosen package from your repository. Commit your source/path/revision configuration and describe any placeholders a new user must replace. Helm users may use an Argo Application rendering a chart or a Flux HelmRelease; Kustomize users may use an Argo Application or Flux Kustomization.

Prove the controller has successfully reconciled a known Git revision. Then change an API replica count directly in the cluster and show it returns to the Git value. Configure drift correction for your chosen route. Do not let two controllers or a second manual Helm installation own the same resources.

## 6. Visible ConfigMap demonstration

Store a background color such as `#e0f2fe` in a ConfigMap. The frontend must actually consume it. Change that color to `#fef3c7` in your configuration, commit and push, then show the browser changing after reconciliation.

No application image rebuild is allowed for this demonstration. Explain whether the implementation rereads a mounted file or needs a controlled Pod restart. Do not confuse a local file edit with a pushed Git change. If you use an environment variable, document how the rollout is triggered; environment variables do not update inside existing processes.

## 7. Operate and recover

Demonstrate a rolling application update using a new image, data persistence across a database Pod replacement, and recovery from one deliberate failure. Acceptable failures include a missing image, incorrect Service selector or broken readiness path. Record the symptom, diagnostic evidence, cause, repair and verification. Restore the correct Git configuration; a temporary live patch alone is not a permanent GitOps repair.

## Submission

Submit repository URL, architecture diagram, exact startup commands, tests and CI run URL, image references, packaging files, GitOps configuration and observed revision, color-change evidence, saved-record persistence evidence, rolling update/drift evidence, incident report and cleanup instructions. Redact credentials. Identify destructive operations, retained storage, and known limitations. A peer must be able to reproduce the deployment on their own kind cluster.
