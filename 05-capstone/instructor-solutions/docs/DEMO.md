# Instructor presentation: 20–30 minutes

Prepare START.md and **one** GITOPS.md route before class. Keep four views available: browser at localhost:8080, terminal for kubectl, the Git configuration file, and Argo UI or Flux status. Use Kustomize + Argo for the shortest first demonstration. Helm and Flux alternatives behave similarly but use the commands noted below.

All shell commands run from instructor-solutions/ on kind-capstone. These are live operations on the disposable demo namespace; do not run against another cluster.

## 1. Show the whole request path — 3 minutes

```bash
kubectl config current-context
kubectl get nodes
kubectl -n capstone get deploy,statefulset,svc,pvc
```

Explain: two frontend Pods, two API Pods, one database Pod. Services connect them. Only the frontend is port-forwarded. Submit “This note should survive a database Pod replacement.” Refresh the page and show it still exists.

Ask: “Is the note in the browser, the API process, or PostgreSQL?”

## 2. Prove database persistence — 3 minutes

```bash
kubectl -n capstone delete pod db-0
kubectl -n capstone rollout status statefulset/db --timeout=180s
kubectl -n capstone rollout status deployment/api --timeout=180s
```

After db-0 is recreated and API readiness recovers, refresh the note list. The record remains because db-data is mounted again. API requests can fail briefly while this single database restarts; two API replicas do not make the database highly available.

Ask: “Would deleting the entire kind cluster preserve this data?” Answer: no; kind stores it in node containers. Do not delete the PVC during this demonstration.

## 3. The main visual demo: change the background — 5 minutes

Show the current ConfigMap:

```bash
kubectl -n capstone get configmap frontend-theme -o yaml
curl -fsS http://localhost:8080/config
```

Open the Git configuration file for the active package:

- Kustomize: deployment/kustomize/overlays/demo/theme.yaml.
- Helm: deployment/helm/campus-board/values.yaml.

Change blue `#e0f2fe` to yellow `#fef3c7` (or use green `#dcfce7` if you already tested yellow). Commit and push only that file as shown in GITOPS.md. Keep the page visible while the controller reconciles.

For Argo, observe Application revision and sync/health in the UI; optionally request a hard refresh. For Flux, request the matching Kustomization or HelmRelease reconciliation. Then run:

```bash
kubectl -n capstone get configmap frontend-theme -o jsonpath='{.data.BACKGROUND_COLOR}'
curl -fsS http://localhost:8080/config
kubectl -n capstone get pods
```

Explain the three stages: Git revision changes, ConfigMap changes, mounted file reaches the frontend. The browser polls the frontend every five seconds. Wait for propagation; do not promise an instantaneous change. Two replicas may briefly see different file revisions, but they converge. The port-forward generally targets one Pod until restarted.

Expected: the background changes and the displayed color label updates. Pod names need not change. No Go compilation, container build or image promotion occurred.

Ask: “Why would the same update not appear automatically if the color were only an environment variable?”

## 4. Demonstrate drift repair — 3 minutes

The desired API replica count is two. Make an intentionally temporary cluster-only change:

```bash
kubectl -n capstone scale deployment/api --replicas=4
kubectl -n capstone get deployment api -w
```

Show the controller returning the desired count to two. Argo uses selfHeal; Flux Kustomization periodically reapplies desired configuration; Flux HelmRelease uses its explicitly enabled drift detection. Stop the watch with Ctrl+C.

Ask: “Where would we change replicas if we wanted four permanently?” Answer: the selected package configuration in Git.

## 5. Demonstrate a rolling application update — 5 minutes

For the local instructor route, build a distinguishable API binary and load it before promoting it:

```bash
docker build -f api/Dockerfile --build-arg VERSION=v2 -t campus-api:v2 api
kind load docker-image campus-api:v2 --name capstone
```

Change only the API image: Kustomize overlay newTag v2, or Helm api.image campus-api:v2. Commit/push that file and watch:

```bash
kubectl -n capstone get pods -w
```

Once the new desired image appears in the Deployment, verify rollout and version:

```bash
kubectl -n capstone rollout status deployment/api --timeout=180s
kubectl -n capstone get replicasets
kubectl -n capstone port-forward service/api 8081:8080
```

In another terminal `curl http://localhost:8081/version` should return v2. Stop the API forward after checking. The browser should still save and list notes during the update; maxUnavailable=0 keeps old ready replicas until replacements are available, subject to actual capacity and dependency health.

For the student publication route, change API source, let CI publish, then promote its actual Docker Hub SHA tag. Use the SHA returned by /version as evidence instead of v2. Never promote a tag before it exists.

## 6. Break one release and recover — 5 minutes

In the chosen Git package, change the API image tag to `missing`. Commit only this mistake and push it. After the controller applies it, run:

```bash
kubectl -n capstone get pods
kubectl -n capstone describe pods -l app=campus-api
```

Expected: new Pods fail to pull; old healthy replicas may continue serving. An Argo Application can be Synced while health is still Progressing or eventually Degraded. Flux may report an unsuccessful health check. Do not delete the database to repair an image problem.

Revert **the specific bad commit** and push. If it is the latest commit with no unrelated work, `git revert HEAD` is appropriate; otherwise use its recorded commit SHA. Request reconciliation if desired. Verify the restored image, rollout status and successful note submission.

Close by asking students to explain the failure evidence, why the DB did not need replacement, and why the repair belongs in Git.

## After class

Keep the cluster for another demo, or follow START.md's explicit data-aware cleanup. Reset the desired color/image in Git if you want the next class to start from blue/v1. Never rewrite shared Git history to hide the failed-release example; the revert is useful teaching evidence.
