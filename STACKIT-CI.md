# STACKIT Git CI setup

The executable workflows live at repository root in `.forgejo/workflows/`.
The `golang-app/` directory is the Git repository root. Push its contents, including the hidden `.forgejo/` directory. Do not put workflows in a plain `workflows/` directory or include a `golang-app/` prefix in workflow paths. Use `main` as the publication branch.
The Go application and container builds use the existing source and Dockerfiles.

## Enable Actions and runners

1. Create or open your repository in STACKIT Git and enable **Repository Actions**.
2. In the STACKIT Portal, open **Developer Platform → Git → your instance → Continuous integration & Runners**, enable runners, and complete the fee-based order. Alternatively, register your own runner with Docker support and change `runs-on` to its label.
3. These workflows select `stackit-docker`. Check that jobs can access a Docker daemon. The capstone also requires Buildx and QEMU/binfmt support for its `linux/amd64,linux/arm64` build; confirm this in the first runner execution.

See [Enable Actions](https://docs.stackit.cloud/products/developer-platform/git/how-tos/enable-actions-in-stackit-git/) and [STACKIT runners](https://docs.stackit.cloud/products/developer-platform/git/getting-started/pipeline-runners/).

## Configure repository Actions secrets

Add these through your STACKIT Git repository's Actions settings. Do not commit credentials.

| Type | Name | Value |
|---|---|---|
| Secret | `DOCKERHUB_USER` | Docker Hub username for all three images |
| Secret | `DOCKERHUB_TOKEN` | Docker Hub token allowed to push that user's `course-app`, `campus-api`, and `campus-frontend` images |

The capstone publishes `docker.io/<DOCKERHUB_USER>/campus-api:sha-<commit>` and `docker.io/<DOCKERHUB_USER>/campus-frontend:sha-<commit>`. The course publishes `docker.io/<DOCKERHUB_USER>/course-app:sha-<commit>`. Both workflows share the same two secrets; no GitHub registry credentials or namespace variable are needed.

Create the `course-app`, `campus-api`, and `campus-frontend` repositories under your Docker Hub account. For the public-image lab, make them public. Use the actual published image names in Helm/Kustomize. Private images require Kubernetes image pull credentials.

The workflows do not rely on GitHub `permissions` declarations: Forgejo does not implement these identically. Scope registry tokens at their issuing service. Checkout uses `persist-credentials: false`. Pull requests run tests and builds; publication is restricted to `main` pushes (and manual capstone runs on `main`).

## Push source

Create the remote repository first. Use SSH authentication or a Git HTTPS credential helper. The helper does not create a remote repository or store credentials.

```bash
# From the golang-app Git repository root, if this folder is not yet a Git repository:
git init -b main
export STACKIT_REPO_URL='https://YOUR-STACKIT-GIT-HOST/OWNER/REPOSITORY.git'
git remote add origin "$STACKIT_REPO_URL"

# Review the files before committing. Include deletions when migrating existing workflows.
git add -A
git diff --cached
git commit -m "Configure STACKIT Git CI"
bash 01-app-ci/scripts/push-source.sh
```

For an existing checkout, retain its Git history and verify `origin` points to your STACKIT repository. The helper refuses a mismatched remote, a branch other than `main` or `kaloyan-app`, or uncommitted files; it pushes without force. On later updates, export `STACKIT_REPO_URL`, commit your changes, and rerun it.

## Verify the first runs

- Open the repository **Actions** tab. Check that both workflows are discovered and that jobs get a `stackit-docker` runner.
- Course CI tests/builds on pull requests and pushes to `main` and `kaloyan-app`; only main pushes publish `DOCKERHUB_USER/course-app:sha-<commit>`.
- Capstone CI tests both services and their Docker builds on every push to `main` or `kaloyan-app`. Pull requests run when API/frontend source or its workflow changes. Only main pushes publish both architectures; manual runs on `main` can also publish.
- Confirm the resulting tags exist in Docker Hub before updating Kubernetes manifests. Confirm a pull request builds without executing publication steps.
- Update complete repository URLs in Argo CD/Flux examples to your STACKIT clone URL. Private Git repositories require controller read credentials; this is separate from registry authentication.

Local YAML/shell validation does not prove managed-runner Docker, QEMU, registry credentials, or publication. Those require a successful run on your STACKIT instance.
