# Start the instructor demonstration

## 1. Prepare an isolated cluster

Use Docker, kind, kubectl, Go, Git, Python 3 and OpenSSL from the earlier modules. Helm 3.17.1 is needed only for the Helm route; Flux CLI 2.5.1 only for Flux. Allow roughly 6 CPUs/10 GB RAM for this cluster with a GitOps controller. If keeping the course cluster running, allow more. Use a free local port 8080.

From the **repository root**:

```bash
cd 05-capstone/instructor-solutions
kind create cluster --name capstone --config deployment/kind-config.yaml
kubectl config use-context kind-capstone
kubectl get nodes -o wide
kubectl get pods -A
kubectl cluster-info
```

Expected nodes: capstone-control-plane, capstone-worker, capstone-worker2. The configuration still creates exactly three nodes; the cluster name controls their prefix. Check `kubectl config current-context` before every session. All remaining commands in this guide run from **instructor-solutions/**. In another terminal, use the same directory/context.

If the cluster already exists, do not recreate it: select kind-capstone and inspect it. The namespace must be dedicated to this demonstration.

## 2. Test and build the two images

```bash
(cd api && go test ./... && go build -o /tmp/campus-api .)
(cd frontend && go test ./... && go build -o /tmp/campus-frontend .)
bash scripts/build-load.sh
```

The scripts are executable: `./scripts/build-load.sh` and `./scripts/create-secret.sh` also work; no sudo is needed. If a downloaded copy loses executable permissions, invoking scripts with bash still works. Finish the build/load and Secret steps before applying the manifests.

The script visibly builds api/Dockerfile and frontend/Dockerfile and loads campus-api:v1 and campus-frontend:v1 into every capstone node. These local images let you test before publishing. PostgreSQL's official image is pulled from Docker Hub when its Pod starts.

```bash
kubectl create namespace capstone --dry-run=client -o yaml | kubectl apply -f -
bash scripts/create-secret.sh
```

The Secret script generates a random password, sends it to Kubernetes via a temporary protected file, removes that file, and retains an existing Secret on reruns. It does not write credentials into the repository. If you keep the PVC, keep the same Secret: changing POSTGRES_PASSWORD on a populated volume does **not** rotate the database password. That variable initializes an empty database. See the [official PostgreSQL image documentation](https://hub.docker.com/_/postgres).

## 3. Choose one package

For the shortest instructor demo, use **Kustomize**:

```bash
kubectl kustomize deployment/kustomize/overlays/demo > /tmp/campus-board.yaml
kubectl apply -k deployment/kustomize/overlays/demo
```

Alternatively, on a fresh namespace with no existing application deployment, use **Helm**:

```bash
helm lint deployment/helm/campus-board
helm template campus-board deployment/helm/campus-board -n capstone > /tmp/campus-board.yaml
helm upgrade --install campus-board deployment/helm/campus-board -n capstone --wait --timeout 5m
mkdir -p /tmp/campus-charts
helm package deployment/helm/campus-board --destination /tmp/campus-charts
```

Both routes create fixed names api, frontend, db, db-data and frontend-theme in capstone. **Do not run both routes against the same namespace.** A second Helm release in that namespace would also collide. To compare routes later, use the explicit data-aware cleanup below and a fresh cluster.

## 4. Wait for the system

```bash
kubectl -n capstone rollout status statefulset/db --timeout=300s
kubectl -n capstone rollout status deployment/api --timeout=300s
kubectl -n capstone rollout status deployment/frontend --timeout=300s
kubectl -n capstone get pods,svc,pvc
kubectl -n capstone get pods -o wide
kubectl -n capstone port-forward service/frontend 8080:8080
```

Keep forwarding running. Open **http://localhost:8080**. Expected: the background color configured by your chosen package, a note field and a list. The Kustomize demo currently uses yellow; Helm defaults to light blue. Add “Hello from the instructor.” In a second terminal run:

```bash
bash scripts/smoke.sh
```

Expected: health success, a JSON background color, a newly created note and `PASS: note persisted and listed`. This writes one real demo note per run. The list returns the latest 100 notes.

## 5. Try the visible color change locally

Before enabling GitOps, change `BACKGROUND_COLOR` in `deployment/kustomize/overlays/demo/theme.yaml` to a different color, for example `#dcfce7` (green). The current demo overlay retains your tested `#fef3c7` yellow setting, then:

```bash
kubectl apply -k deployment/kustomize/overlays/demo
kubectl -n capstone get configmap frontend-theme -o yaml
```

For the Helm route instead change `backgroundColor` in deployment/helm/campus-board/values.yaml and run:

```bash
helm upgrade campus-board deployment/helm/campus-board -n capstone --wait --timeout 5m
```

Watch the page change to the configured color. Allow up to about two minutes for the mounted ConfigMap to propagate, plus the five-second browser poll. You can query `curl http://localhost:8080/config` to separate browser behavior from server configuration. Do not use subPath or convert this setting to an environment variable if you want this live reload behavior.

Now follow [GITOPS.md](GITOPS.md), then present [DEMO.md](DEMO.md). Once GitOps is active, make permanent changes in Git instead of running direct apply/upgrade commands.

## Troubleshooting

| Symptom | Check and repair |
|---|---|
| ImagePullBackOff | `kubectl -n capstone describe pods`; check local images loaded into capstone or public Docker Hub references. |
| API waiting for database | `kubectl -n capstone logs deployment/api`; check db readiness, db-credentials and Service DNS. |
| Database permission/startup failure | `kubectl -n capstone logs db-0`; inspect PVC, node disk space and security context. Use the pinned Debian image with UID 999. |
| PVC Pending | `kubectl -n capstone describe pvc db-data`; standard uses WaitForFirstConsumer, so wait for db scheduling; inspect local-path-storage Pods. |
| Password failures after recreating Secret | Restore the matching credential or deliberately reset disposable DB data; changing the Secret alone does not alter a saved PostgreSQL password. |
| Frontend works but notes fail | `kubectl -n capstone get endpointslices`; inspect API readiness and logs; db failure removes API from ready Service endpoints. |
| Background remains old | Inspect ConfigMap, then /config; wait for volume propagation and check the active package's color file. |
| Port-forward disconnects after a rollout | Restart the forwarding command: it may be attached to a Pod that was replaced. |

## Cleanup: choose whether to keep data

First stop forwarding and remove GitOps management using GITOPS.md. To preserve a backup before any destructive cleanup:

```bash
kubectl -n capstone exec db-0 -- pg_dump -U campus campus > /tmp/campus-board-backup.sql
```

The local demo uses the database container's local socket for this command. Protect the dump if it contains non-demo data.

To remove workload resources **while keeping the PVC and Secret**, after removing GitOps ownership:

```bash
# Kustomize route only: explicitly list deletions so db-data is retained.
kubectl -n capstone delete deployment api frontend
kubectl -n capstone delete statefulset db
kubectl -n capstone delete service api frontend db
kubectl -n capstone delete configmap frontend-theme
```

For the manual Helm route, use `helm uninstall campus-board -n capstone`; the chart marks db-data with `helm.sh/resource-policy: keep`. Retained PVC ownership can complicate later release changes, so inspect it before adopting it. Do not use `kubectl delete -k`: an explicit delete can remove PVCs despite controller pruning annotations.

For a **complete disposable reset**, after saving anything needed:

```bash
kind delete cluster --name capstone
```

This destroys all capstone workloads, Secrets and node-local database storage. It leaves the course cluster and Docker Hub images alone. For a fresh cluster, build/load images and generate a fresh Secret again.
