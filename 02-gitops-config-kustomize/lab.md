# Demo: the same app, different replica counts

**Time:** 20 minutes. **Starting state:** Module 01 image loaded into the three-node kind cluster. No GitOps controller should manage these namespaces yet. Run from `golang-app/`.

## 1. Preview

**What you will do:** inspect what the dev overlay produces.

**Instructions:**

```bash
kubectl config use-context kind-course
kubectl kustomize 02-gitops-config-kustomize/overlays/dev
```

**Write down — Q1.** Find the rendered Deployment. What are its namespace and replica count?

**Done when:** you can point to those two fields in the output.

## 2. Deploy both environments

```bash
kubectl apply -k 02-gitops-config-kustomize/overlays/dev
kubectl apply -k 02-gitops-config-kustomize/overlays/prod
kubectl -n course-dev rollout status deployment/go-app --timeout=120s
kubectl -n course-prod rollout status deployment/go-app --timeout=120s
kubectl -n course-dev get deployment go-app
kubectl -n course-prod get deployment go-app
```

**Write down — Q2.** How many replicas are Ready in each environment? Did you need to copy the Deployment file?

**Done when:** dev shows 1/1 and prod shows 2/2.

## 3. Make one small change

Open `02-gitops-config-kustomize/overlays/dev/kustomization.yaml`. Change only `count: 1` to `count: 2`.

```bash
kubectl apply -k 02-gitops-config-kustomize/overlays/dev
kubectl -n course-dev rollout status deployment/go-app --timeout=120s
kubectl -n course-dev get deployment go-app
kubectl -n course-prod get deployment go-app
```

**Write down — Q3.** Which environment changed, and why?

**Done when:** dev now has two Ready replicas. Restore its count to 1 and apply dev again so Module 03 starts from the expected configuration.

## Quick troubleshooting

If a rollout times out, run `kubectl -n course-dev describe pods`. Check that course-app:v1 was loaded with `kind load docker-image course-app:v1 --name course`. A valid render cannot guarantee an image exists.

Optional two-minute failure: back up the dev overlay, change `../../base` to `../../missing`, and render it. Observe the error, then restore the path. Rendering does not change existing Pods.

## Cleanup and submission

Keep dev and prod for the GitOps modules. If stopping the course here, `kubectl delete namespace course-dev course-prod` removes both environments and everything in those namespaces.

Submit Q1–Q3 and the before/after Deployment output. This demo does not require image promotion or rollout rollback; those are in the extended lab.
