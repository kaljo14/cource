# Student lab — 00-lab-setup

## Prerequisites and goal

Read [theory](README.md) and complete the preceding module. Work from `golang-app/`, context `kind-course`. Build, inspect, break and repair the system. Keep a terminal available for port-forward and use a second for curl. Ctrl+C stops watches and forwarding.

## Part 1 — Create and inspect

**What you will do:** create and inspect a kubernets cluster.

**Instructions:**

```bash
bash 00-lab-setup/scripts/check.sh
kind create cluster --name course --config 00-lab-setup/kind-config.yaml
kubectl config use-context kind-course
kubectl get nodes -o wide
kubectl get pods -A
kubectl cluster-info
```

**Write down — Q1.** Which components run on the control plane?

**Done when:** Three Ready nodes named course-control-plane, course-worker, course-worker2.

## Part 2 — Namespaces and failure

**What you will do:** create a namespace and a pod on the namespace

**Instructions:**

```bash
kubectl create namespace inspection
kubectl -n inspection run probe --image=busybox:1.37 -- sleep 3600
kubectl -n inspection get pods -o wide
kubectl -n default get pod probe
kubectl -n inspection describe pod probe
```

**Write down — Q2.** Why does the default namespace lookup fail? Find the assigned worker and owner references.

**Done when:** The Pod runs on a worker the wrong-namespace command returns NotFound.

## Cleanup

```bash
kubectl delete namespace inspection
```

## Submission

Submit answers to Q1–Q2, commands, one before/after failure observation, and a short explanation of the repair. Redact credentials.
