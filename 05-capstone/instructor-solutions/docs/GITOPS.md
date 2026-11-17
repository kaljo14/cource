# Add GitOps: choose exactly one route

Start with START.md completed, working application and capstone context. Commands run from instructor-solutions/. This guide supports four alternatives; apply **only one** application integration. Keep the namespace and db-credentials outside controller ownership.

## 1. Put the solution in your repository

Push the contents of golang-app/ (the Git repository root) to your own STACKIT Git repository, branch main. If this workspace is not yet a Git repository, initialize and commit it from its root before adding your STACKIT Git remote. The supplied paths and root workflow use the complete course layout. You can also copy this solution into its own repository; see the standalone-path table below.

Replace the complete GitHub placeholder URL in the GitOps manifest you choose, and in flux-source.yaml if using Flux, with your STACKIT Git clone URL. For private repositories, configure read credentials in Argo CD or Flux. Commit and push. No public URL is preconfigured because it must be your own repository.

For a quick instructor demonstration, keep campus-api:v1/campus-frontend:v1 in Git; the images loaded into all capstone nodes work with IfNotPresent. **For the student final submission**, configure DOCKERHUB_USER and DOCKERHUB_TOKEN using [STACKIT CI setup](../../../STACKIT-CI.md), then run the root Capstone application CI workflow on main and make both Docker Hub repositories public. Update the deployment configuration to the resulting full tags:

- Kustomize: set each images entry's newName to `docker.io/<DOCKERHUB_USER>/campus-api` or campus-frontend, and newTag to `sha-<successful-CI-commit>`.
- Helm: set api.image and frontend.image in deployment/helm/campus-board/values.yaml to the complete published references.

Commit the promotion only after the image build succeeds. A color-only configuration commit does not need a new image build. The workflow's manual dispatch on main is available if no source path changed. Do not attach registry credentials to this public-image demo.


### Same repository or a standalone repository

The commands in these guides run from the solution directory. GitOps paths, however, are relative to the **remote Git repository root**.

| Setting | This course repository | Solution copied to its own repository root |
|---|---|---|
| Argo Kustomize source.path | `05-capstone/instructor-solutions/deployment/kustomize/overlays/demo` | `deployment/kustomize/overlays/demo` |
| Flux Kustomization spec.path | `./05-capstone/instructor-solutions/deployment/kustomize/overlays/demo` | `./deployment/kustomize/overlays/demo` |
| Argo Helm source.path | `05-capstone/instructor-solutions/deployment/helm/campus-board` | `deployment/helm/campus-board` |
| Flux Helm chart.spec.chart | `./05-capstone/instructor-solutions/deployment/helm/campus-board` | `./deployment/helm/campus-board` |

For a standalone repository, also copy the course-root `.forgejo/workflows/capstone-ci.yaml` into its root `.forgejo/workflows/` and remove the `05-capstone/instructor-solutions/` prefix from its paths, contexts and working directories. The two Docker build contexts become api/ and frontend/. Update the repository URL in the chosen Argo manifest or Flux source. START.md's initial cd into the course subfolder can then be skipped; deployment/kind-config.yaml is included locally.

### Updating an existing installation after this refactor

Keep the existing namespace, Secret, PVC and workloads. Commit and push the moved files first. Until its source path is updated, a running controller may report that the old path is missing. Use **only the command matching the integration you already run**, from instructor-solutions/ on kind-capstone:

```bash
# Argo + Kustomize
kubectl -n argocd patch application campus-board --type=merge -p '{"spec":{"source":{"path":"05-capstone/instructor-solutions/deployment/kustomize/overlays/demo"}}}'

# Argo + Helm
kubectl -n argocd patch application campus-board --type=merge -p '{"spec":{"source":{"path":"05-capstone/instructor-solutions/deployment/helm/campus-board"}}}'

# Flux + Kustomize
kubectl -n flux-system patch kustomization campus-board --type=merge -p '{"spec":{"path":"./05-capstone/instructor-solutions/deployment/kustomize/overlays/demo"}}'

# Flux + Helm
kubectl -n flux-system patch helmrelease campus-board --type=merge -p '{"spec":{"chart":{"spec":{"chart":"./05-capstone/instructor-solutions/deployment/helm/campus-board"}}}}'
```

These patches change only the path, preserving your live repository URL and other settings. The matching files in deployment/gitops/ already contain the new paths; keep your actual repository URL in those files too. If another GitOps resource manages these controller objects, change its Git configuration as well. Inspect status and wait for the pushed revision. Do not delete application data to resolve a path error.

## 2A. Argo CD with Kustomize — recommended first demo

The initial Kustomize resources can be reconciled by Argo directly. There must be no Helm release or Flux controller managing them.

```bash
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.14.3/manifests/install.yaml
kubectl -n argocd rollout status deployment/argocd-server --timeout=300s
kubectl apply -f deployment/gitops/argo-kustomize.yaml
kubectl -n argocd get application campus-board -w
```

Expect Synced and Healthy after the correct repository revision is read and workloads are ready. Stop the watch with Ctrl+C. For the dashboard:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | python3 -c 'import sys,base64; print(base64.b64decode(sys.stdin.read()).decode())'
kubectl -n argocd port-forward service/argocd-server 8443:443
```

Open https://localhost:8443, accept the self-signed certificate only for this known local endpoint, and log in as admin. Keep the password private.

## 2B. Argo CD with Helm

For the cleanest route, start a fresh capstone cluster, build/load images and create the namespace/Secret as in START.md, but **skip manual Helm installation**. Then install Argo as above and apply:

```bash
kubectl apply -f deployment/gitops/argo-helm.yaml
kubectl -n argocd get application campus-board -w
```

Argo renders the chart in Git and manages the resulting resources. This is not a Helm CLI release; do not expect helm list to show it and do not use helm upgrade against it. To follow this route after testing a manual Helm release, prefer a fresh disposable cluster. Avoid silently adopting an existing Helm installation and its retained PVC.

## 2C. Flux with Kustomize

There must be no Argo Application or HelmRelease managing capstone. Install Flux in this cluster:

```bash
flux check --pre
flux install --version=v2.5.1
flux check
kubectl apply -f deployment/gitops/flux-source.yaml
kubectl apply -f deployment/gitops/flux-kustomize.yaml
flux reconcile kustomization campus-board --with-source
flux get sources git
flux get kustomizations
```

The existing manually applied Kustomize resources can be reconciled. Look for Ready=True and the intended revision. For errors, inspect `kubectl -n flux-system describe kustomization campus-board` and `flux logs --level=error`.

## 2D. Flux with Helm

Start a fresh prepared cluster as in route 2B and skip manual chart installation. Install Flux as in 2C, then:

```bash
kubectl apply -f deployment/gitops/flux-source.yaml
kubectl apply -f deployment/gitops/flux-helm.yaml
flux reconcile source git campus-board
flux reconcile helmrelease campus-board --with-source
flux get helmreleases
```

The HelmRelease lives in flux-system and installs into capstone. It references a chart in the GitRepository. `reconcileStrategy: Revision` picks up chart content changes at new Git revisions for this teaching workflow. `driftDetection.mode: enabled` enables repair of live resource changes; merely using a HelmRelease is not sufficient to demonstrate automatic drift correction without that configuration. Do not also create flux-kustomize.yaml.

## 3. Change the color through Git

Edit the file for your chosen package:

- Kustomize: `deployment/kustomize/overlays/demo/theme.yaml`, BACKGROUND_COLOR.
- Helm: `deployment/helm/campus-board/values.yaml`, backgroundColor.

Change to `#fef3c7`. Stage **only the chosen file**, commit and push, for example:

```bash
git add deployment/kustomize/overlays/demo/theme.yaml
git commit -m "Make campus board background pale yellow"
git push
```

For Helm use `git add deployment/helm/campus-board/values.yaml` instead. Do not apply the resources manually while demonstrating GitOps. Observe controller status, the ConfigMap and the browser. You can request a fresh Argo comparison with:

```bash
kubectl -n argocd annotate application campus-board argocd.argoproj.io/refresh=hard --overwrite
```

For Flux use the matching reconcile command from your route. Wait for the Git revision to be reconciled, then for the projected file and browser poll. The color can lag a healthy controller status briefly.

## 4. Stop management without deleting database data

Do this before switching tools or performing manual cleanup. Do not switch tools during the main class demo.

Argo route: these Applications intentionally have no deletion finalizer, so deleting the Application leaves workloads running:

```bash
kubectl -n argocd delete application campus-board
```

Flux Kustomize route: disable pruning before deleting the reconciliation object:

```bash
kubectl -n flux-system patch kustomization campus-board --type=merge -p '{"spec":{"prune":false}}'
kubectl -n flux-system delete kustomization campus-board
kubectl -n flux-system delete gitrepository campus-board
```

Flux Helm route: deleting the HelmRelease normally **uninstalls the release**. The PVC has a keep policy but the application workload is removed. If you only want to stop changes temporarily, use:

```bash
flux suspend helmrelease campus-board
```

Suspension retains the management object; do not introduce another controller over the same release. For final removal use `kubectl -n flux-system delete helmrelease campus-board`, wait for it to finish, then delete its GitRepository. Inspect the retained PVC before proceeding with START.md cleanup.

No pruning annotation protects storage from an explicit namespace deletion or deletion of the entire kind cluster. Export data first if needed.
