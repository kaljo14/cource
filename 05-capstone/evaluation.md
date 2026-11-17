# Evaluation — 100 points

| Area | Points | Observable criteria |
|---|---:|---|
| Three-tier functionality | 20 | Frontend saves/lists records through API (8); real DB persistence (6); input validation and meaningful tests (6). |
| Containers and CI | 15 | Reproducible separate images (5); tests/builds in Actions (5); trusted-only publication and traceable Docker Hub references (5). |
| Kubernetes and packaging | 15 | Own three-node cluster and correct internal Services (5); coherent Helm chart or Kustomize base/overlay (5); reproducible deployment instructions (5). |
| Reliability and persistence | 15 | Probes and requests/limits (5); two API/frontend replicas and rolling update (5); PVC persistence demonstrated and limitations explained (5). |
| GitOps | 15 | Correct source/revision/path and healthy reconciliation (6); live replica drift repaired (5); single ownership and permanent Git repair (4). |
| Configuration and security | 10 | ConfigMap color change visibly reaches browser without image rebuild (5); non-root app containers, Secret handling and localhost access (5). |
| Explanation and recovery | 10 | Incident evidence and successful repair (5); architecture, Q1–Q6, cleanup/data-loss explanation (5). |

Award partial credit for independently verified subcriteria. Either packaging tool and either GitOps controller can earn full credit; implementing all combinations gives no additional points. A copied implementation without a correct explanation earns implementation credit only where the student can demonstrate it. Screenshots complement, but do not replace, reproducible commands.
