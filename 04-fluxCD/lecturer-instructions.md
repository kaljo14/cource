# Lecturer instructions: Flux from a private repository and custom branch

Task 04 fetches the private STACKIT Git repository's `kaloyan-app` branch and deploys Task 02's production overlay. Commands run from the cloned repository root (`golang-app/`), unless stated otherwise. Use Bash for the interactive prompts below.

The two Flux resources have separate responsibilities:

| File | Responsibility |
|---|---|
| `04-fluxCD/manifests/source.yaml` | Repository URL, branch `kaloyan-app`, and authentication Secret |
| `04-fluxCD/manifests/kustomization.yaml` | Builds `./02-gitops-config-kustomize/overlays/prod` from that source |

The Task 02 overlay includes `../../base` and deploys into `course-prod`. Do not add a `golang-app/` prefix to the path. Flux runs in the cluster and requires network access to your Git host; credentials on your laptop are not inherited.

## 1. Prepare the branch and image

Confirm that `origin` points to your STACKIT repository. Switch to `kaloyan-app` if it exists locally; if it exists only remotely, fetch it and create a tracking branch:

```bash
git fetch origin
git switch --track origin/kaloyan-app
```

If the branch does not exist anywhere, create it from your intended starting commit with `git switch -c kaloyan-app`. If it already exists locally, use `git switch kaloyan-app` instead.

Edit `04-fluxCD/manifests/source.yaml`: replace the placeholder URL with the repository's HTTPS clone URL. Leave `ref.branch: kaloyan-app` and `secretRef.name: course-git-auth`.

For the registry-based demo, set the image in `02-gitops-config-kustomize/overlays/prod/kustomization.yaml` to an image that has already been published:

```yaml
images:
- name: course-app
  newName: docker.io/YOUR_DOCKERHUB_USER/course-app
  newTag: sha-ACTUAL_SUCCESSFUL_BUILD_COMMIT
```

CI tests/builds both `main` and `kaloyan-app`, but publishes images only from `main`. A configuration-only commit on `kaloyan-app` can reuse an image published from main. Do not set the image tag to a new custom-branch commit unless that image actually exists.

For the original local kind exercise, you can instead retain `course-app:v1` after loading that image into all course cluster nodes. Private Docker Hub images require a separate image pull Secret in `course-prod`; Git credentials do not authenticate container pulls.

## 2. Install Flux in the course cluster

Complete Task 03's documented Argo handoff before proceeding; its Application must no longer manage the course resources. Select and verify the cluster, then install Flux if it is not already installed:

```bash
kubectl config use-context kind-course
kubectl config current-context
flux check --pre
flux install --version=v2.5.1
flux check
```

This uses the course's pinned Flux version. If Flux is already installed, run `flux check` rather than reinstalling or downgrading it. The `flux-system` namespace must exist before the next step.

## 3. Create private Git credentials

Create a STACKIT Git access token for an account with read access to this repository. Restrict it to repository read access where supported; this exercise does not need Git write permission. Use this Git token, not the Docker Hub token or a STACKIT cloud API credential.

Run locally in Bash:

```bash
read -r -p 'STACKIT HTTPS clone URL: ' STACKIT_GIT_URL
read -r -p 'STACKIT Git username: ' STACKIT_GIT_USER
read -r -s -p 'STACKIT Git access token: ' STACKIT_GIT_TOKEN
printf '\n'

flux create secret git course-git-auth \
  --namespace=flux-system \
  --url="$STACKIT_GIT_URL" \
  --username="$STACKIT_GIT_USER" \
  --password="$STACKIT_GIT_TOKEN"

unset STACKIT_GIT_TOKEN
```

The URL must match `spec.url` in `source.yaml`. The command creates the credentials directly in Kubernetes. Do not commit an unencrypted Secret or token. The Secret and GitRepository both live in `flux-system`.

Flux supports HTTPS username/password credentials in a referenced Secret, and its CLI accepts the token as the password. See [GitRepository authentication](https://fluxcd.io/flux/components/source/gitrepositories/#secret-reference) and [flux create secret git](https://fluxcd.io/flux/cmd/flux_create_secret_git/).

## 4. Commit the configuration and connect Flux

Review your edits, then commit and push them to the tracked branch:

```bash
git add 04-fluxCD/manifests/source.yaml \
  02-gitops-config-kustomize/overlays/prod/kustomization.yaml
git diff --cached
git commit -m "Configure private Git source and production image"
git push -u origin kaloyan-app

kubectl apply -f 04-fluxCD/manifests/source.yaml
kubectl apply -f 04-fluxCD/manifests/kustomization.yaml
flux reconcile kustomization course --with-source
flux get sources git
flux get kustomizations
kubectl -n course-prod rollout status deployment/go-app --timeout=180s
```

Expected: both Flux resources report `Ready=True`; the source revision identifies `kaloyan-app` and its fetched commit. Healthy Pods require the selected image to exist and be accessible to cluster nodes.

## 5. Demonstrate Git-driven reconciliation

Change `MESSAGE=Hello prod` in the Task 02 production overlay to `MESSAGE=Hello from kaloyan-app`. Commit that file and push to `kaloyan-app`. Then run:

```bash
flux reconcile kustomization course --with-source
flux get kustomizations
kubectl -n course-prod rollout status deployment/go-app --timeout=180s
```

Flux rebuilds the overlay; Kustomize's generated ConfigMap reference changes, causing the Deployment to roll out. Editing a local file without pushing does not affect Flux. A commit made only to main does not change this branch's desired state.

## 6. Use a different custom branch

Push the desired configuration to that branch first. Change `spec.ref.branch` in `04-fluxCD/manifests/source.yaml`, commit it, and apply the source file again:

```bash
kubectl apply -f 04-fluxCD/manifests/source.yaml
flux reconcile kustomization course --with-source
```

The Git Secret and Task 02 path stay the same. This lab does not have Flux managing its own Task 04 manifests: changes to the source URL, branch, or Flux path need `kubectl apply` to update the cluster. The pipeline's branch allowlist is separate; a third branch needs a workflow change only if you want CI to run on its pushes.

## Troubleshooting

| Symptom | Check |
|---|---|
| GitRepository not Ready | Exact clone URL, remote branch exists, token read access, Secret name/namespace, cluster network access |
| Certificate error | Git server certificate trust; configure the required CA instead of disabling verification |
| Source Ready, Kustomization not Ready | Task 02 path, YAML validity, health-check errors |
| ImagePullBackOff | Published image tag, Docker Hub visibility and image pull credentials |
| Changes not appearing | Commit pushed to tracked branch; source revision and last applied revision |

```bash
kubectl -n flux-system describe gitrepository course
kubectl -n flux-system describe kustomization course
flux logs --level=error
```

Do not print Secret contents in lecture recordings. Follow the student lab's cleanup when finished; disabling pruning before removing Flux resources retains the application workloads.
