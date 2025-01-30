.PHONY: frontend
frontend:
	@cd frontend && npm run dev

.PHONY: backend
backend:
	@cd backend/cmd/api && go run .

.PHONY: run
run: frontend backend
