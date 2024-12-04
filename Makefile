.PHONY: frontend/run
frontend/run:
	@cd frontend && bun run dev & open http://localhost:5173

.PHONY: backend/run
backend/run:
	@cd backend/cmd/api && go run .