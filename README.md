# RoutineTask

A simple task manager (routine + ad-hoc tasks), built as a hands-on DevOps learning project.

## Stack
- Backend: Go (chi router, SQLite via modernc.org/sqlite)
- Frontend: Next.js (TypeScript, Tailwind)
- Containers: Docker

## Run locally

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

## Run with Docker

```bash
cd backend && docker build -t routinetask-backend .
cd ../frontend && docker build -t routinetask-frontend .
```

## Roadmap
- [x] Go API + Next.js UI
- [x] Dockerfiles (multi-stage)
- [ ] Docker Compose
- [ ] CI/CD (GitLab CI)
- [ ] Kubernetes + Helm
- [ ] ArgoCD (GitOps)
