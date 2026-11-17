# Frontend

A standalone Go module with no external dependencies. static/index.html defines the page, static/styles.css controls its responsive layout, and static/app.js handles notes and the live ConfigMap theme. frontend.go serves the page, proxies /api/ to the API and reads the mounted background-color file. main.go reads API_URL and COLOR_FILE; server.go handles graceful shutdown.

From this directory: `go test ./...` and `docker build -t campus-frontend:v1 .`.

Defaults: API_URL=http://api:8080, COLOR_FILE=/config/BACKGROUND_COLOR; listen port 8080. See [startup instructions](../docs/START.md) for the complete demo.

## Show the updated interface in kind

Static assets are embedded in the Go binary, so changing HTML/CSS/JavaScript requires a new frontend image. Run from instructor-solutions/:

```bash
docker build -t campus-frontend:v2 frontend
kind load docker-image campus-frontend:v2 --name capstone
```

For Kustomize, change only campus-frontend's newTag to v2 in deployment/kustomize/overlays/demo/kustomization.yaml. For Helm, set frontend.image to campus-frontend:v2 in deployment/helm/campus-board/values.yaml. If GitOps is active, commit and push the chosen configuration and wait for reconciliation. Otherwise apply the overlay or upgrade the existing chart as in START.md.

After the new image appears in the Deployment:

```bash
kubectl -n capstone rollout status deployment/frontend --timeout=180s
kubectl -n capstone port-forward service/frontend 8080:8080
```

Stop any old forwarding process before starting this one, then reload http://localhost:8080. Use a hard refresh if the browser cached old assets. With GHCR, promote the real published SHA tag instead of the local v2 tag. A background-color-only ConfigMap change still requires no image rebuild.
