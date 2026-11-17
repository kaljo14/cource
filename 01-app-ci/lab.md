# Student lab — 01-app-ci

## Prerequisites and goal

Read [theory](README.md) and complete the preceding module. Work from `golang-app/`, context `kind-course`. Build, inspect, break and repair the system. Keep a terminal available for port-forward and use a second for curl. Ctrl+C stops watches and forwarding.

## Part 1 — Run and test

**What you will do:** run and test.

**Instructions:**

```bash
(cd 01-app-ci/app && go test -race ./... && go build -o /tmp/course-server .)
MESSAGE="Hello local" /tmp/course-server
```

**Write down — Q1.** In a second terminal curl localhost:8080/ and /readyz. What does each report? Stop with Ctrl+C.

**Done when:** JSON message and ready response; process exits on SIGINT.

## Part 2 — Container and deployment

**What you will do:** container and deployment.

**Instructions:**

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

**Write down — Q2.** Why is kind load needed even though docker images lists the image?

**Done when:** Two Ready Pods, ideally spread across workers; localhost:8080 responds.

## Part 3 — Scale, update and diagnose

**What you will do:** scale, update and diagnose.

**Instructions:**

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

**Write down — Q3.** Which event explains the failure? Why can old replicas still serve requests? Repair to v2.

**Done when:** Image pull failure observed and repaired; rollout completes.

## Part 4 — Publish CI

**What you will do:** publish ci.

**Instructions:**

```bash
# At repository root, follow instructor-commands.md to configure STACKIT Git
# and the DOCKERHUB_USER / DOCKERHUB_TOKEN repository secrets.
bash 01-app-ci/scripts/push-source.sh
```

**Write down — Q4.** What tag and permissions did the workflow use? Why must it live at repository root?

**Done when:** Tests and publication succeed; Docker Hub image exists for the commit.

## Cleanup

These commands remove only the named lab resources; retain state explicitly noted for later modules. Namespace deletion removes every object inside it.

```bash
kubectl delete namespace course-app
# Keep local images for subsequent modules.
```

## Submission

Submit answers to Q1–Q4, commands, one before/after failure observation, and a short explanation of the repair. Redact credentials.
