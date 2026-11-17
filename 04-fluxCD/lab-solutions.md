# Instructor solutions — 04-fluxCD

Outputs described below are representative, not captured cluster results. Discover actual Pod names, IPs, revisions and passwords; never submit credentials.

## Part 1

**Q1.** It prevents accidental overlapping reconciliation. Existing dev objects remain; Flux targets separate course-prod resources.

Expected observation: No Argo Application owns course resources; Flux controllers ready.

Commands and verification:

```bash
kubectl -n argocd delete application course --ignore-not-found
kubectl -n argocd get applications
flux check --pre
flux install --version=v2.5.1
flux check
```

## Part 2

**Q2.** Git fetch can succeed while rendering, applying, or health checks fail. Each resource reports a different stage.

Expected observation: Both report Ready=True at the pushed revision.

Commands and verification:

```bash
# Replace repository URL in source.yaml. Update prod overlay to your public Docker Hub
# image and successful sha tag as in Module 03; commit and push.
kubectl apply -f 04-fluxCD/manifests/source.yaml
kubectl apply -f 04-fluxCD/manifests/kustomization.yaml
flux reconcile source git course
flux reconcile kustomization course --with-source
flux get sources git
flux get kustomizations
kubectl -n course-prod rollout status deployment/go-app
```

## Part 3

**Q3.** Existing workloads keep running because the build failed before apply. kubectl apply -f 04-fluxCD/manifests/kustomization.yaml then flux reconcile kustomization course --with-source restores the path.

Expected observation: Failed reconciliation diagnosed, repaired and Ready again.

Commands and verification:

```bash
kubectl -n course-prod scale deployment/go-app --replicas=4
flux reconcile kustomization course
# Commit a prod MESSAGE change and push; reconcile again.
kubectl -n flux-system patch kustomization course --type=merge -p '{"spec":{"path":"./missing"}}'
flux reconcile kustomization course
flux logs --level=error
kubectl -n flux-system describe kustomization course
```

## Troubleshooting and cleanup
Check context, namespace, events, Pod image and controller conditions in that order. Logs require a started container; use describe for scheduling/pull failures. `kubectl logs --previous` helps after crashes. A valid alternative is inspecting the same conditions in the UI, accompanied by CLI evidence.

```bash
kubectl -n flux-system patch kustomization course --type=merge -p '{"spec":{"prune":false}}'
kubectl -n flux-system delete kustomization course
kubectl -n flux-system delete gitrepository course
# Workloads retained. Delete course-dev/course-prod only when finished; this removes all contents.
```
