# From Go Application to Kubernetes: CI/CD, GitOps, and Platform Engineering

A hands-on university course for students with basic Linux, Git and programming knowledge. Difficulty: guided intermediate tasks with one manageable failure at a time. Approximately nine and a half teaching hours plus capstone preparation and optional exercises.

## Learning path
- [00-lab-setup: Lab setup](00-lab-setup/README.md) — 60 minutes
- [01-app-ci: Application, containers and CI](01-app-ci/README.md) — 150 minutes
- [02-gitops-config-kustomize: Configuration with Kustomize](02-gitops-config-kustomize/README.md) — 30 minutes
- [03-argoCD: GitOps with Argo CD](03-argoCD/README.md) — 90 minutes
- [04-fluxCD: GitOps with Flux](04-fluxCD/README.md) — 90 minutes
- [04a-helm: Package and deploy with Helm](04a-helm/README.md) — 30 minutes
- [05: Three-tier capstone](05-capstone/README.md) — 120-minute kickoff/demo + 4–6 hours independent work

Read theory, attempt the lab, collect evidence, then consult solutions. Run shell commands from this directory. Preserve the enclosing repository layout: workflows live in the parent .forgejo/workflows and GitOps paths start with golang-app/. Create a STACKIT Git repository containing this entire workspace; main is the assumed branch. See [STACKIT CI setup](../STACKIT-CI.md) for runners and registry credentials. Managed STACKIT runners are fee-based; local kind labs can run without them.

## Start
```bash
cd golang-app
bash 00-lab-setup/scripts/check.sh
kind create cluster --name course --config 00-lab-setup/kind-config.yaml
kubectl config use-context kind-course
kubectl get nodes -o wide
kubectl get pods -A
kubectl cluster-info
```

See the [curriculum](curriculum.md) and the capstone [instructor demonstration](05-capstone/instructor-solutions/docs/DEMO.md). The original CNI/CSI reference material remains unchanged in ../example/. These materials follow its inspect–explain–break–repair pattern without requiring its unrelated modules.

## Final cleanup
`bash 00-lab-setup/scripts/delete.sh` destroys the entire course cluster, all workloads and node-local data. Save submitted evidence first. It does not delete GitHub repositories or GHCR packages.
