# Student lab — 03-argoCD

Estimated duration: 45–100 minutes within the module allocation.

## Prerequisites and goal
Read [theory](README.md) and complete the preceding module. Work from `golang-app/`, context `kind-course`. Build, inspect, break and repair the system. Keep a terminal available for port-forward and use a second for curl. Ctrl+C stops watches and forwarding.

## Part 1 — Install and access

**What you will do:** install and access.

**Instructions:**

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.14.3/manifests/install.yaml
kubectl -n argocd rollout status deployment/argocd-server --timeout=300s
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | python3 -c "import sys,base64; print(base64.b64decode(sys.stdin.read()).decode())"
kubectl -n argocd port-forward service/argocd-server 8443:443
```

**Write down — Q1.** Open https://localhost:8443 and log in as admin. Why is the certificate untrusted?

**Done when:** UI accessible only through localhost; initial password kept private.

## Part 2 — Connect Git

**What you will do:** connect git.

**Instructions:**

```bash
# Edit application.yaml: replace REPLACE_OWNER/REPLACE_REPO with your public repository.
# In dev overlay set newName to ghcr.io/<lowercase-owner>/<lowercase-repo>/go-app
# and newTag to sha-<full successful CI commit SHA>. Commit and push these changes.
kubectl apply -f 03-argoCD/manifests/application.yaml
kubectl -n argocd get applications -w
```

**Write down — Q2.** What is the difference between Synced and Healthy?

**Done when:** Application reaches Synced and Healthy; dev serves published version.

## Part 3 — Reconcile and break

**What you will do:** reconcile and break.

**Instructions:**

```bash
kubectl -n course-dev scale deployment/go-app --replicas=4
kubectl -n course-dev get deployment/go-app -w
# Change dev MESSAGE in Git, commit and push; observe a new ConfigMap and rollout.
# Then commit an invalid image tag; inspect Application status and Pod events.
# Revert only that last commit using git revert HEAD and push.
```

**Write down — Q3.** Which changes persist and why?

**Done when:** Live scale is reverted; committed message persists; bad image is repaired by Git revert.

## Cleanup
These commands remove only the named lab resources; retain state explicitly noted for later modules. Namespace deletion removes every object inside it.

```bash
kubectl -n argocd delete application course
# No deletion finalizer: workloads remain. Do not delete course-dev here.
```

## Submission
Submit answers to Q1–Q3, commands, one before/after failure observation, and a short explanation of the repair. Redact credentials.
