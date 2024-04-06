.PHONY: help
help:
	@echo "Usage: make [ TARGETS ]"
	@echo ""
	@echo "Targets:"
	@echo "  build     		- Build the application, swagger document will be generated"
#	@echo "  gen-swagger 		- Generate swagger document"
	@echo "  run			- Run the frontend and backend"
	@echo "  deploy		-deploy the application on my site"

.PHONY:	build
build:
	cd frontend
	vite build
	npm install

.PHONY: start-backend start-frontend start

start-backend:
	@cd scripts && ./check_port.sh
	@cd src/backend && go run main.go & echo "backend started"

start-frontend:
	@cd src/frontend && npm install && npm run dev & echo "frontend started"

run: start-backend start-frontend
	@echo "All services started"

.PHONY: deploy
deploy:
	@cd scripts && ./deploy.sh
	@echo "Application deployed"