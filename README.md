# RoutineTask

A simple task manager (routine + ad-hoc tasks), built as a hands-on, step-by-step DevOps learning project — from a basic Go + Next.js app all the way to Kubernetes and GitOps.

## Stack

| Layer | Technology | Notes |
|---|---|---|
| Backend | Go | `chi` router, `database/sql` (no ORM), `modernc.org/sqlite` (pure-Go driver, no CGO) |
| Frontend | Next.js | TypeScript, Tailwind, App Router, `output: standalone` |
| Database | SQLite | File-based, persisted via volume/PVC |
| Containers | Docker | Multi-stage builds for both services |
| Orchestration (local) | Docker Compose | Healthchecks, named volume, configurable env vars |
| Orchestration (cluster) | Kubernetes (Minikube) | Deployments, Services, ConfigMap, PVC, probes |
| CI/CD | GitLab CI (self-hosted) | Test → build → push to Container Registry |
| Source control | GitLab (self-hosted) + GitHub (mirror) | Branch → PR → merge workflow |

## Architecture

```
┌─────────────┐      ┌─────────────┐
│   Frontend  │─────▶│   Backend   │─────▶ SQLite (PVC)
│  (Next.js)  │      │    (Go)     │
└─────────────┘      └─────────────┘
   NodePort              ClusterIP
   (30080)                (8080)
```

In Kubernetes, both run inside the `routinetask` namespace. The backend is only reachable from inside the cluster (`ClusterIP`); the frontend is exposed via `NodePort` for local access.

## Run locally (no containers)

Backend:
```bash
cd backend
go run main.go
```

Frontend:
```bash
cd frontend
npm install
npm run dev
```

Open http://localhost:3000

## Run with Docker Compose

```bash
docker compose up --build -d
```

- Frontend: http://localhost:3000
- Backend health: http://localhost:8080/health

Stop:
```bash
docker compose down
```

Stop and wipe the database too:
```bash
docker compose down -v
```

## Run on Kubernetes (Minikube)

### Prerequisites
- `kubectl` and `minikube` installed
- Docker images built and available (either pulled from the GitLab Container Registry or loaded directly via `minikube image load`, see **Known limitations** below)

### Apply manifests

```bash
cd k8s
kubectl apply -f 00-namespace.yaml
kubectl apply -f 01-backend-configmap.yaml
kubectl apply -f 02-backend-pvc.yaml
kubectl apply -f 03-backend-deployment.yaml
kubectl apply -f 04-backend-service.yaml
kubectl apply -f 05-frontend-deployment.yaml
kubectl apply -f 06-frontend-service.yaml
```

Or apply everything in one go (filenames are numbered so `kubectl` applies them in the right order):
```bash
kubectl apply -f k8s/
```

### What each manifest does

| File | Kind | Purpose |
|---|---|---|
| `00-namespace.yaml` | Namespace | Isolates all project resources under `routinetask` |
| `01-backend-configmap.yaml` | ConfigMap | Non-sensitive backend config (`DB_PATH`, `PORT`, `CORS_ALLOWED_ORIGIN`), injected as env vars |
| `02-backend-pvc.yaml` | PersistentVolumeClaim | 1Gi `ReadWriteOnce` volume for the SQLite file, so data survives Pod restarts |
| `03-backend-deployment.yaml` | Deployment | Runs the backend container; mounts the PVC at `/data`; has liveness/readiness probes on `/health` |
| `04-backend-service.yaml` | Service (`ClusterIP`) | Internal-only DNS name (`backend`) so the frontend can reach the backend inside the cluster |
| `05-frontend-deployment.yaml` | Deployment | Runs the frontend container; liveness/readiness probes on `/` |
| `06-frontend-service.yaml` | Service (`NodePort`) | Exposes the frontend outside the cluster on port `30080` |

### Accessing the app

Minikube runs as its own Docker container, so `localhost` from your host machine doesn't reach the cluster's NodePort directly. Use port-forwarding:

```bash
# terminal 1
kubectl port-forward -n routinetask svc/frontend 30080:3000

# terminal 2
kubectl port-forward -n routinetask svc/backend 8080:8080
```

Then open http://localhost:30080

(Port `30080` is used deliberately on both sides so it matches `CORS_ALLOWED_ORIGIN` in the ConfigMap.)

### Useful commands

```bash
kubectl get pods -n routinetask
kubectl get svc -n routinetask
kubectl logs -n routinetask -l app=backend --tail=50
kubectl logs -n routinetask -l app=frontend --tail=50
kubectl describe pod -n routinetask -l app=backend
```

## CI/CD (GitLab)

Pipeline (`.gitlab-ci.yml`) runs on every push to `main`:

1. **test stage** (parallel)
   - `backend-test`: `go vet` + `go build`
   - `frontend-test`: `npm ci` + `npm run lint` + `npm run build`
2. **build stage** (parallel, `main` only)
   - `backend-docker-build`: builds and pushes `backend` image to the project's Container Registry
   - `frontend-docker-build`: builds and pushes `frontend` image, tagged with both `latest` and the short commit SHA

Requires a self-hosted GitLab Runner (Docker executor, socket-mounted — see **Known limitations**).

## Known limitations / decisions

These are real constraints hit while building this (mostly due to restricted access to some registries from Iran) and the workarounds chosen — documented here the way a real project would log architectural trade-offs:

- **Ingress skipped.** The `ingress-nginx` controller's images live on `registry.k8s.io`, which wasn't reachable. Mirrored the images via `m.daocloud.io`, but the controller also needs digest-pinned images for its admission webhook Jobs, which added more friction than it was worth for a learning lab. Using `NodePort` + `kubectl port-forward` instead. Revisit if a reliable internal registry mirror becomes available.
- **Images loaded directly into Minikube**, not pulled from the registry *inside* the cluster. The GitLab Container Registry's JWT auth flow requires reaching the GitLab instance's external URL (`http://localhost:8929`), which isn't resolvable from inside Minikube's network namespace. Workaround: `docker pull` on the host → `docker tag` → `minikube image load`. In a real cluster (or with a public GitLab URL), this wouldn't be an issue.
- **`imagePullPolicy: IfNotPresent` set explicitly** on both Deployments. The `latest` tag defaults to `imagePullPolicy: Always` in Kubernetes, which would otherwise force a registry pull every time — defeating the point of `minikube image load`.
- **`HOSTNAME=0.0.0.0` set explicitly** on the frontend container. Next.js `standalone` mode reads the `HOSTNAME` env var to decide what interface to bind to; Kubernetes auto-injects `HOSTNAME` as the Pod's name, which made the server unreachable via `port-forward` (it was only listening on the Pod's own hostname, not `localhost`/`0.0.0.0`).
- **GitLab Runner uses Docker socket mounting**, not true Docker-in-Docker. Simpler to set up and sufficient for a lab; a production setup would typically use DinD with `--privileged` for better isolation.
- **GOPROXY overridden** to `https://goproxy.io` in CI (and recommended locally) because `proxy.golang.org` isn't reachable.
- **Container Registry helper image overridden** to `gitlab/gitlab-runner-helper` (Docker Hub) instead of the default `registry.gitlab.com/...` image, for the same reason.

## Roadmap

- [x] Go backend + Next.js frontend
- [x] Dockerize (multi-stage builds)
- [x] Docker Compose (healthcheck, persistent volume, configurable env)
- [x] GitLab self-hosted + Runner + CI/CD pipeline (test → build → push)
- [x] Kubernetes (Minikube) — Namespace, ConfigMap, PVC, Deployments, Services, probes
- [ ] Helm (package the manifests above into a reusable chart)
- [ ] ArgoCD (GitOps — sync deployments automatically from the Git repo)
