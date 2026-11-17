# Student lab — 02-gitops-config-kustomize

Estimated duration: 45–100 minutes within the module allocation.

## Prerequisites and goal
Read [theory](README.md) and complete the preceding module. Work from `golang-app/`, context `kind-course`. Build, inspect, break and repair the system. Keep a terminal available for port-forward and use a second for curl. Ctrl+C stops watches and forwarding.

## Part 1 — Render and deploy

**What you will do:** render and deploy.

**Instructions:**

```bash
kubectl create namespace course-base
kubectl -n course-base apply -k 02-gitops-config-kustomize/base
kubectl kustomize 02-gitops-config-kustomize/overlays/dev
kubectl apply -k 02-gitops-config-kustomize/overlays/dev
kubectl apply -k 02-gitops-config-kustomize/overlays/prod
kubectl -n course-dev get deploy,cm
kubectl -n course-prod get deploy,cm
```

**Write down — Q1.** Why are replica counts and ConfigMap names different?

**Done when:** Dev has one replica; prod has two; messages differ.

## Part 2 — Promote and roll back

**What you will do:** promote and roll back.

**Instructions:**

```bash
cp 02-gitops-config-kustomize/overlays/dev/kustomization.yaml /tmp/dev-overlay.yaml
# Edit dev newTag to v2, render, apply, and inspect rollout.
kubectl apply -k 02-gitops-config-kustomize/overlays/dev
kubectl -n course-dev rollout status deployment/go-app
kubectl -n course-dev rollout undo deployment/go-app
cp /tmp/dev-overlay.yaml 02-gitops-config-kustomize/overlays/dev/kustomization.yaml
kubectl apply -k 02-gitops-config-kustomize/overlays/dev
```

**Write down — Q2.** Why must the file also be restored after rollout undo?

**Done when:** Live and file state both return to v1.

## Part 3 — Break rendering

**What you will do:** break rendering.

**Instructions:**

```bash
cp 02-gitops-config-kustomize/overlays/dev/kustomization.yaml /tmp/dev-overlay.yaml
# Change ../../base to ../../missing, then run:
kubectl kustomize 02-gitops-config-kustomize/overlays/dev
```

**Write down — Q3.** Does a render failure alter the cluster? Repair using the backup.

**Done when:** Renderer reports missing resource; restored overlay renders.

## Cleanup
These commands remove only the named lab resources; retain state explicitly noted for later modules. Namespace deletion removes every object inside it.

```bash
kubectl delete namespace course-base
# Keep dev and prod for GitOps.
```

## Submission
Submit answers to Q1–Q3, commands, one before/after failure observation, and a short explanation of the repair. Redact credentials.
