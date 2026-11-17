# Instructor solutions — 03-argoCD

Outputs described below are representative, not captured cluster results. Discover actual Pod names, IPs, revisions and passwords; never submit credentials.

## Part 1

**Q1.** The default installation uses a self-signed certificate. Accept only for this known local lab endpoint.

Expected observation: UI accessible only through localhost; initial password kept private.

Commands and verification:

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.14.3/manifests/install.yaml
kubectl -n argocd rollout status deployment/argocd-server --timeout=300s
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | python3 -c "import sys,base64; print(base64.b64decode(sys.stdin.read()).decode())"
kubectl -n argocd port-forward service/argocd-server 8443:443
```

## Part 2

**Q2.** Synced describes configuration equality, Healthy describes resource operation. Inspect status conditions and kubectl -n course-dev describe pods if health fails.

Expected observation: Application reaches Synced and Healthy; dev serves published version.

Commands and verification:

```bash
# Edit application.yaml: replace REPLACE_OWNER/REPLACE_REPO with your public repository.
# In dev overlay set newName to docker.io/<dockerhub-user>/course-app
# and newTag to sha-<full successful CI commit SHA>. Commit and push these changes.
kubectl apply -f 03-argoCD/manifests/application.yaml
kubectl -n argocd get applications -w
```

## Part 3

**Q3.** Self-heal restores replica count from Git. A committed change becomes desired state. A bad tag can be Synced but Degraded/Progressing; revert and push restores the prior image.

Expected observation: Live scale is reverted; committed message persists; bad image is repaired by Git revert.

Commands and verification:

```bash
kubectl -n course-dev scale deployment/go-app --replicas=4
kubectl -n course-dev get deployment/go-app -w
# Change dev MESSAGE in Git, commit and push; observe a new ConfigMap and rollout.
# Then commit an invalid image tag; inspect Application status and Pod events.
# Revert only that last commit using git revert HEAD and push.
```

## Troubleshooting and cleanup
Check context, namespace, events, Pod image and controller conditions in that order. Logs require a started container; use describe for scheduling/pull failures. `kubectl logs --previous` helps after crashes. A valid alternative is inspecting the same conditions in the UI, accompanied by CLI evidence.

```bash
kubectl -n argocd delete application course
# No deletion finalizer: workloads remain. Do not delete course-dev here.
```
