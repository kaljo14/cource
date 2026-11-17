# Instructor solutions — 02-gitops-config-kustomize

Outputs described below are representative, not captured cluster results. Discover actual Pod names, IPs, revisions and passwords; never submit credentials.

## Part 1

**Q1.** Overlay replicas override base; generated content hashes reflect different MESSAGE values.

Expected observation: Dev has one replica; prod has two; messages differ.

Commands and verification:

```bash
kubectl create namespace course-base
kubectl -n course-base apply -k 02-gitops-config-kustomize/base
kubectl kustomize 02-gitops-config-kustomize/overlays/dev
kubectl apply -k 02-gitops-config-kustomize/overlays/dev
kubectl apply -k 02-gitops-config-kustomize/overlays/prod
kubectl -n course-dev get deploy,cm
kubectl -n course-prod get deploy,cm
```

## Part 2

**Q2.** Undo changes live Pod template only. A later apply would restore v2 unless configuration also changes.

Expected observation: Live and file state both return to v1.

Commands and verification:

```bash
cp 02-gitops-config-kustomize/overlays/dev/kustomization.yaml /tmp/dev-overlay.yaml
# Edit dev newTag to v2, render, apply, and inspect rollout.
kubectl apply -k 02-gitops-config-kustomize/overlays/dev
kubectl -n course-dev rollout status deployment/go-app
kubectl -n course-dev rollout undo deployment/go-app
cp /tmp/dev-overlay.yaml 02-gitops-config-kustomize/overlays/dev/kustomization.yaml
kubectl apply -k 02-gitops-config-kustomize/overlays/dev
```

## Part 3

**Q3.** No API mutation occurs while rendering. cp /tmp/dev-overlay.yaml 02-gitops-config-kustomize/overlays/dev/kustomization.yaml repairs it.

Expected observation: Renderer reports missing resource; restored overlay renders.

Commands and verification:

```bash
cp 02-gitops-config-kustomize/overlays/dev/kustomization.yaml /tmp/dev-overlay.yaml
# Change ../../base to ../../missing, then run:
kubectl kustomize 02-gitops-config-kustomize/overlays/dev
```

## Troubleshooting and cleanup
Check context, namespace, events, Pod image and controller conditions in that order. Logs require a started container; use describe for scheduling/pull failures. `kubectl logs --previous` helps after crashes. A valid alternative is inspecting the same conditions in the UI, accompanied by CLI evidence.

```bash
kubectl delete namespace course-base
# Keep dev and prod for GitOps.
```
