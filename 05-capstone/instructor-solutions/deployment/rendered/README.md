# Saved render

`demo.yaml` is your existing rendered snapshot, preserved unchanged during the refactor. It may differ from the current source overlay. Use `deployment/kustomize/overlays/demo` as the source of truth; do not point GitOps at this directory. Generate fresh previews in /tmp when following the guides.
