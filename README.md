# From Go Application to Kubernetes: CI/CD, GitOps, and Platform Engineering

A hands-on university course for students with basic Linux, Git and programming knowledge. Difficulty: guided intermediate tasks with one manageable failure at a time. Approximately nine teaching hours plus capstone preparation and optional exercises.

## Learning path

- [00-lab-setup: Lab setup](00-lab-setup/README.md)
- [01-app-ci: Application, containers and CI](01-app-ci/README.md)
- [02-gitops-config-kustomize: Configuration with Kustomize](02-gitops-config-kustomize/README.md)
- [03-argoCD: GitOps with Argo CD](03-argoCD/README.md)
- [04-fluxCD: GitOps with Flux](04-fluxCD/README.md)
- [05: Three-tier capstone](05-capstone/README.md)

Read theory, attempt the lab, collect evidence, then consult solutions. Run shell commands from this directory. This directory (`golang-app/`) is the Git repository root. Push its contents, including `.forgejo/workflows/`; workflow and GitOps paths are relative to this directory. CI runs on pushes to `main` and `kaloyan-app`, with publication from `main` only. See [STACKIT CI setup](STACKIT-CI.md) for runners and registry credentials. Managed STACKIT runners are fee-based; local kind labs can run without them.

## Start

```bash
# Run from the cloned repository root.
bash 00-lab-setup/scripts/check.sh
kind create cluster --name course --config 00-lab-setup/kind-config.yaml
kubectl config use-context kind-course
kubectl get nodes -o wide
kubectl get pods -A
kubectl cluster-info
```

See the [curriculum](curriculum.md) and the capstone [instructor demonstration](05-capstone/instructor-solutions/docs/DEMO.md). The course is self-contained; the original reference material in the parent workspace is not required.

## Final cleanup

`bash 00-lab-setup/scripts/delete.sh` destroys the entire course cluster, all workloads and node-local data. Save submitted evidence first. It does not delete STACKIT Git repositories or Docker Hub images.
