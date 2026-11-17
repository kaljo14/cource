```bash
# Run from the STACKITUNI repository root.

# Start kind.
kind get clusters | grep -qx course || kind create cluster --name course --config golang-app/00-lab-setup/kind-config.yaml
kubectl config use-context kind-course
kubectl wait --for=condition=Ready nodes --all --timeout=180s
kubectl create namespace course-app --dry-run=client -o yaml | kubectl apply -f -

# Run the pipeline steps locally.
bash golang-app/01-app-ci/scripts/build-app.sh
bash golang-app/01-app-ci/scripts/push-dockerhub.sh

# Load the local image, build manifests, and deploy.
kind load docker-image course-app:v1 --name course
kubectl kustomize golang-app/01-app-ci/k8s > /tmp/course-app.yaml
kubectl -n course-app apply -f /tmp/course-app.yaml
kubectl -n course-app rollout status deployment/go-app --timeout=180s
kubectl -n course-app port-forward service/go-app 8080:80

# Run in a second terminal.
curl http://localhost:8080/
curl http://localhost:8080/readyz
```

```bash
# Configure STACKIT Actions, runners, and secrets as described in STACKIT-CI.md.
# Create an empty repository in STACKIT Git and copy its clone URL.
export STACKIT_REPO_URL='https://YOUR-STACKIT-GIT-HOST/OWNER/REPOSITORY.git'

# For a new workspace only:
# git init -b main
# git remote add origin "$STACKIT_REPO_URL"
# Review and commit changes before publication.
git add -A
git diff --cached
git commit -m "Configure STACKIT Git CI"
bash golang-app/01-app-ci/scripts/push-source.sh
# Open the repository Actions tab and wait for successful publication.

# Pull the CI image and load it into kind.
read -r -p 'Docker Hub username: ' DOCKERHUB_USER
read -r -s -p 'Docker Hub token (read): ' DOCKERHUB_TOKEN
printf '\n'
printf '%s' "$DOCKERHUB_TOKEN" | docker login --username "$DOCKERHUB_USER" --password-stdin
unset DOCKERHUB_TOKEN
IMAGE_TAG="sha-$(git rev-parse HEAD)"
REGISTRY_IMAGE="$DOCKERHUB_USER/course-app"
docker pull "$REGISTRY_IMAGE:$IMAGE_TAG"
kind load docker-image "$REGISTRY_IMAGE:$IMAGE_TAG" --name course

# Build manifests with the CI image and deploy.
DEPLOY_DIR=$(mktemp -d)
cp golang-app/01-app-ci/k8s/*.yaml "$DEPLOY_DIR/"
cat >> "$DEPLOY_DIR/kustomization.yaml" <<YAML
images:
  - name: course-app
    newName: $REGISTRY_IMAGE
    newTag: $IMAGE_TAG
YAML
kubectl kustomize "$DEPLOY_DIR" > /tmp/course-app-ci.yaml
kubectl --context kind-course -n course-app apply -f /tmp/course-app-ci.yaml
kubectl --context kind-course -n course-app rollout status deployment/go-app --timeout=180s
kubectl --context kind-course -n course-app port-forward service/go-app 8080:80
```
