# Short demo solutions

Representative results, not captured cluster output.

**Q1.** Namespace course-dev, replicas 1. `kubectl kustomize` prints the merged configuration without an API write.

**Q2.** Dev has 1/1 Ready replicas; prod has 2/2 once rollouts finish. Both reuse base/deployment.yaml. The overlays override replicas and namespace.

**Q3.** Only dev changes from one to two replicas. Its overlay is independent of prod. Changing the shared base could affect both, depending on their overrides.

Restore dev's `count: 1`, then verify:

```bash
kubectl apply -k 02-gitops-config-kustomize/overlays/dev
kubectl -n course-dev rollout status deployment/go-app --timeout=120s
kubectl -n course-dev get deployment go-app
kubectl -n course-prod get deployment go-app
```

Expect dev 1/1 and prod 2/2. If not ready, inspect `kubectl -n course-dev describe pods` and events. Image errors require a valid loaded image; scheduling failures require sufficient capacity. Do not repair by editing unrelated namespaces.

For the optional failure, a missing-base error is a local rendering failure. Restore ../../base and rerender; no rollout is needed to undo a failed preview.

Retain both namespaces for later modules. At the end of the course, deleting course-dev/course-prod removes their entire contents. The optional [extended solutions](extended-lab-solutions.md) cover image rollback and the longer exercise.
