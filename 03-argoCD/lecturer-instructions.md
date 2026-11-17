# Lecturer instructions: Argo CD setup with a private Git repository

Task 03 deploys Task 02's **dev overlay** from a custom branch. The example tracks `kaloyan-app` and creates application resources in `course-dev`. Task 04 uses Flux and the prod overlay instead.

Run commands from the cloned repository root (`golang-app/`). You need the course kind cluster, kubectl, Git, and access to the remote repository. The browser-based setup below does not require the Argo CD CLI.

## 1. Understand the configuration

`03-argoCD/manifests/application.yaml` connects three settings:

| Setting | Meaning |
|---|---|
| `spec.source.repoURL` | Exact HTTPS clone URL of your repository |
| `spec.source.targetRevision` | Branch to deploy: `kaloyan-app` |
| `spec.source.path` | `02-gitops-config-kustomize/overlays/dev` within that branch |

The dev overlay includes `../../base`, sets one replica, and supplies its environment configuration. The path has no `golang-app/` prefix because that directory is already the Git repository root.

The Application lives in `argocd`; its workloads live in `course-dev`. Argo CD reads Git and applies Kubernetes resources. It does not build or publish the container image.

## 2. Prepare the repository and image

Use the repository that actually contains your course files. For example, if you use the GitHub repository from your recent push, the URL is `https://github.com/kaljo14/cource.git`. For STACKIT Git, copy its HTTPS clone URL. Credentials must come from the chosen Git host, not the other one.

Fetch and switch to the custom branch:

```bash
git fetch origin
# If the branch exists locally:
git switch kaloyan-app
```

If it exists only on the remote, use `git switch --track origin/kaloyan-app` instead. If it does not exist yet, create it from your intended starting commit with `git switch -c kaloyan-app`.

Edit the Application's `repoURL` to that clone URL. Leave `targetRevision: kaloyan-app`. Do not embed credentials in the URL.

For the Docker Hub demo, edit `02-gitops-config-kustomize/overlays/dev/kustomization.yaml`:

```yaml
images:
- name: course-app
  newName: docker.io/YOUR_DOCKERHUB_USER/course-app
  newTag: sha-ACTUAL_SUCCESSFUL_BUILD_COMMIT
```

Use a tag that already exists. The STACKIT CI workflows build both main and kaloyan-app, but publish only from main. A custom-branch configuration change can reuse an image previously published from main. Hosting the files on GitHub alone does not run the STACKIT `.forgejo/workflows/` pipelines; publish the image through your configured build environment first.

Alternatively, the local kind lab can keep `course-app:v1` after loading that image into all course nodes. Public Docker Hub images need no pull Secret; private images need one in `course-dev`, referenced by the workload. Git credentials and container registry credentials are separate.

Review and publish the configuration:

```bash
git add 03-argoCD/manifests/application.yaml \
  02-gitops-config-kustomize/overlays/dev/kustomization.yaml
git diff --cached
git commit -m "Configure Argo CD custom branch and dev image"
git push -u origin kaloyan-app
```

## 3. Install Argo CD

Verify that you are targeting the lab cluster:

```bash
kubectl config use-context kind-course
kubectl config current-context
kubectl get nodes
```

For a fresh lab installation:

```bash
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.14.3/manifests/install.yaml
kubectl -n argocd rollout status deployment/argocd-server --timeout=300s
kubectl -n argocd rollout status deployment/argocd-repo-server --timeout=300s
kubectl -n argocd rollout status statefulset/argocd-application-controller --timeout=300s
kubectl -n argocd get pods
```

This retains the course's historical v2.14.3 pin. If Argo CD is already installed, inspect it instead of reinstalling or downgrading it. The destination is the same cluster, so no external cluster registration is needed.

See the [official installation and access guide](https://argo-cd.readthedocs.io/en/stable/getting_started/).

## 4. Open the dashboard

In one terminal, keep this running:

```bash
kubectl -n argocd port-forward service/argocd-server 8443:443
```

Open `https://localhost:8443`. The default lab certificate is self-signed; accept it only for this known local endpoint. In a second terminal, retrieve the initial admin password privately:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath='{.data.password}' | \
  python3 -c 'import sys,base64; print(base64.b64decode(sys.stdin.read()).decode())'
```

Sign in as `admin`. Do not display the password in lecture recordings. If the password was already changed, use the current one. Change the initial password in the account settings, then delete the initial-password Secret after verifying your new login:

```bash
kubectl -n argocd delete secret argocd-initial-admin-secret
```

## 5. Connect the private repository

Create a Git-host access token with permission to read this repository. An account that cannot read the repository cannot grant Argo CD access to it. For this exercise Argo CD needs no write permission. Use a GitHub token for GitHub, or a STACKIT Git token for STACKIT Git; a Docker Hub token cannot authenticate Git.

In the Argo CD dashboard:

1. Open **Settings → Repositories → Connect Repo**.
2. Select **HTTPS** and repository type **git**.
3. Enter the exact clone URL used in the Application.
4. Enter the Git account username and use the token as the password.
5. Connect and confirm the connection status is **Successful**.

Argo CD stores repository credentials as Kubernetes Secrets in its namespace and matches them to the repository URL. Unlike the Flux example, this Application does not need a `secretRef` field. Do not save the token in `temp-token`, a YAML file, or any other tracked file. Your laptop's Git login is not automatically available to Argo CD.

See [private repository authentication](https://argo-cd.readthedocs.io/en/stable/user-guide/private-repositories/) and [declarative repository configuration](https://argo-cd.readthedocs.io/en/stable/operator-manual/declarative-setup/#repositories).

## 6. Create the Application and verify the deployment

```bash
kubectl apply -f 03-argoCD/manifests/application.yaml
kubectl -n argocd get application course -w
```

Open **Applications → course** in the dashboard. Automatic synchronization is enabled in the manifest. Check the selected branch, commit, resource tree, sync status, and health. Stop the watch with Ctrl+C.

```bash
kubectl -n argocd get application course \
  -o jsonpath='{.status.sync.status}{" | "}{.status.health.status}{" | "}{.status.sync.revision}{"\n"}'
kubectl -n course-dev rollout status deployment/go-app --timeout=180s
kubectl -n course-dev get pods,svc
```

Expected: **Synced** means live configuration matches the rendered Git configuration; **Healthy** means the managed resources meet their health checks. A bad image tag can be Synced while Pods fail to start.

To view the app, use another terminal:

```bash
kubectl -n course-dev port-forward service/go-app 8080:80
```

Open `http://localhost:8080`.

## 7. Demonstrate Git changes and drift repair

Change the dev overlay's `MESSAGE=Hello dev` to `MESSAGE=Hello from kaloyan-app`. Commit that file and push to kaloyan-app. Click **Refresh** in the Application dashboard, or request a refresh:

```bash
kubectl -n argocd annotate application course \
  argocd.argoproj.io/refresh=hard --overwrite
```

Observe the revision and rollout. Kustomize changes the generated ConfigMap reference, which updates the Deployment's Pod template. Changes only on main are not deployed by this Application.

Next, demonstrate live drift:

```bash
kubectl -n course-dev scale deployment/go-app --replicas=4
kubectl -n course-dev get deployment/go-app -w
```

With `selfHeal: true`, Argo CD restores the one replica declared in the dev overlay. With `prune: true`, resources removed from this application's desired configuration can be deleted during synchronization. Make lasting changes in Git.

## 8. Track another branch or change the configuration

Push the new branch with its manifests first. Edit `spec.source.targetRevision` in the Application, commit the change, and apply the Application again:

```bash
kubectl apply -f 03-argoCD/manifests/application.yaml
```

The same procedure applies to changing its URL or overlay path. This lab's Application renders only Task 02's overlay; it does not manage its own Task 03 Application file. Committing an Application edit alone does not update the live Application.

Choosing a branch for Argo CD does not change the CI trigger branches or publication policy. Also, do not point Argo and Flux at the same resources simultaneously. The supplied examples separate Argo's course-dev from Flux's course-prod.

## 9. Troubleshoot and hand off to Task 04

| Symptom | Check |
|---|---|
| Repository connection failed | Clone URL, token host and permissions, expiry, cluster network access |
| Revision cannot be resolved | Branch exists on the remote and spelling matches targetRevision |
| Manifest generation failed | Path relative to repository root and valid Task 02 Kustomize files |
| ImagePullBackOff | Published tag, Docker Hub visibility, image pull credentials |
| Synced but not Healthy | Pod events, readiness checks, image or scheduling failures |
| Changes do not appear | Commit pushed to kaloyan-app, Application revision, refresh |

```bash
kubectl -n argocd describe application course
kubectl -n argocd logs deployment/argocd-repo-server --tail=100
kubectl -n course-dev describe pods
```

For the course handoff, delete the supplied Application:

```bash
kubectl -n argocd delete application course
```

The supplied manifest has no resources deletion finalizer, so this direct kubectl deletion leaves its workloads running. Do not select cascading deletion in the UI if you want to retain them. Continue with the [Flux lecturer guide](../04-fluxCD/lecturer-instructions.md); leave namespace deletion until the final lab cleanup.
