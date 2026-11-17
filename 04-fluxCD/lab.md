# Student lab — 04-fluxCD

Estimated duration: 45–100 minutes within the module allocation.

## Prerequisites and goal
Read [theory](README.md) and complete the preceding module. Work from `golang-app/`, context `kind-course`. Build, inspect, break and repair the system. Keep a terminal available for port-forward and use a second for curl. Ctrl+C stops watches and forwarding.

## Part 1 — Safe handoff and installation

**What you will do:** safe handoff and installation.

**Instructions:**

```bash
kubectl -n argocd delete application course --ignore-not-found
kubectl -n argocd get applications
flux check --pre
flux install --version=v2.5.1
flux check
```

**Write down — Q1.** Why remove the Argo Application before further experimentation?

**Done when:** No Argo Application owns course resources; Flux controllers ready.

## Part 2 — Source and reconciliation

Complete the [private-repository setup](lecturer-instructions.md) first.

**What you will do:** source and reconciliation.

**Instructions:**

```bash
# Follow lecturer-instructions.md: set the private repo URL and create course-git-auth.
# Flux tracks kaloyan-app. Use an existing published Docker Hub image SHA in the
# prod overlay; commit and push to kaloyan-app. CI publishes images only from main.
kubectl apply -f 04-fluxCD/manifests/source.yaml
kubectl apply -f 04-fluxCD/manifests/kustomization.yaml
flux reconcile source git course
flux reconcile kustomization course --with-source
flux get sources git
flux get kustomizations
kubectl -n course-prod rollout status deployment/go-app
```

**Write down — Q2.** Why inspect both the source and Kustomization?

**Done when:** Both report Ready=True at the pushed revision.

## Part 3 — Drift and failed reconciliation

**What you will do:** drift and failed reconciliation.

**Instructions:**

```bash
kubectl -n course-prod scale deployment/go-app --replicas=4
flux reconcile kustomization course
# Commit a prod MESSAGE change and push; reconcile again.
kubectl -n flux-system patch kustomization course --type=merge -p '{"spec":{"path":"./missing"}}'
flux reconcile kustomization course
flux logs --level=error
kubectl -n flux-system describe kustomization course
```

**Write down — Q3.** Find the path error and restore the manifest. What happens to running Pods during this failure?

**Done when:** Failed reconciliation diagnosed, repaired and Ready again.

## Cleanup
These commands remove only the named lab resources; retain state explicitly noted for later modules. Namespace deletion removes every object inside it.

```bash
kubectl -n flux-system patch kustomization course --type=merge -p '{"spec":{"prune":false}}'
kubectl -n flux-system delete kustomization course
kubectl -n flux-system delete gitrepository course
# Workloads retained. Delete course-dev/course-prod only when finished; this removes all contents.
```

## Submission
Submit answers to Q1–Q3, commands, one before/after failure observation, and a short explanation of the repair. Redact credentials.
