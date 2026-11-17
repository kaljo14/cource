# Instructor solutions — 01-app-ci

Outputs described below are representative, not captured cluster results. Discover actual Pod names, IPs, revisions and passwords; never submit credentials.

## Part 1

**Q1.** Root returns configured application data; readiness returns 200 when accepting traffic.

Expected observation: JSON message and ready response; process exits on SIGINT.

Commands and verification:

```bash
(cd 01-app-ci/app && go test -race ./... && go build -o /tmp/course-server .)
MESSAGE="Hello local" /tmp/course-server
```

## Part 2

**Q2.** kind nodes have separate containerd image stores. Loading copies the host image into all nodes.

Expected observation: Two Ready Pods, ideally spread across workers; localhost:8080 responds.

Commands and verification:

```bash
docker build --build-arg VERSION=v1 -t course-app:v1 01-app-ci/app
docker run --rm -p 8080:8080 course-app:v1
# After testing with curl in another terminal, stop the container with Ctrl+C.
kind load docker-image course-app:v1 --name course
kubectl create namespace course-app
kubectl -n course-app apply -k 01-app-ci/k8s
kubectl -n course-app rollout status deployment/go-app --timeout=120s
kubectl -n course-app get deploy,rs,pods,svc -o wide
kubectl -n course-app port-forward service/go-app 8080:80
```

## Part 3

**Q3.** ErrImagePull then ImagePullBackOff indicate missing image. maxUnavailable=0 preserves available old Pods. Repair: kubectl -n course-app set image deployment/go-app app=course-app:v2; then rollout status.

Expected observation: Image pull failure observed and repaired; rollout completes.

Commands and verification:

```bash
kubectl -n course-app scale deployment/go-app --replicas=3
docker build --build-arg VERSION=v2 -t course-app:v2 01-app-ci/app
kind load docker-image course-app:v2 --name course
kubectl -n course-app set image deployment/go-app app=course-app:v2
kubectl -n course-app rollout status deployment/go-app
kubectl -n course-app set image deployment/go-app app=course-app:missing
kubectl -n course-app get pods
kubectl -n course-app describe pods -l app=go-app
kubectl -n course-app get events --sort-by=.metadata.creationTimestamp
```

## Part 4

**Q4.** The tag is sha- plus the full commit SHA. The workflow checks out source without persisting Git credentials and uses DOCKERHUB_USER / DOCKERHUB_TOKEN repository secrets to publish to Docker Hub. STACKIT Git discovers workflows under root .forgejo/workflows. Use the repository containing golang-app. Pull requests run the build script; main pushes also run the Docker Hub publish script.

Expected observation: Tests and publication succeed; Docker Hub image exists for the commit.

Commands and verification:

```bash
# Configure repository secrets using instructor-commands.md first.
bash golang-app/01-app-ci/scripts/push-source.sh
```

## Troubleshooting and cleanup
Check context, namespace, events, Pod image and controller conditions in that order. Logs require a started container; use describe for scheduling/pull failures. `kubectl logs --previous` helps after crashes. A valid alternative is inspecting the same conditions in the UI, accompanied by CLI evidence.

```bash
kubectl delete namespace course-app
# Keep local images for subsequent modules.
```
