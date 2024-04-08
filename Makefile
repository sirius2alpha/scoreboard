.PHONY: help
help:
	@echo "Usage: make [ TARGETS ]"
	@echo ""
	@echo "Targets:"
	@echo "  build     		- Build the application, swagger document will be generated"
	@echo "  debug		- Run the backend in debug mode"
	@echo "  deploy		-deploy the application on my site"
#	@echo "  gen-swagger 		- Generate swagger document"
	@echo "  run			- Run the frontend and backend"

.PHONY:	build
build:
	cd app/frontend 
	npm install

.PHONY: debug
debug:
	@cd app/backend
	dlv debug ./app/backend/main.go

.PHONY: deploy
deploy:
	@cd app/frontend/ && npm run build
	scp -r dist root@sirius1y.top:/var/www/scoreboard/frontend/
	@cd app/backend/ && go build -o main
	scp main root@sirius1y.top:/var/www/scoreboard/backend/
	@echo "Application deployed"

.PHONY: run
run:
	@cd scripts && ./check_port.sh
	@cd app/frontend && npm run dev & 
	@echo "frontend started"
	@cd app/backend
	go run main.go &
	@echo "backend started"
	@echo "All services started"