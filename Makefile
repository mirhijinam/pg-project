up: 
	docker-compose up --build
.PHONY: up

down: 
	docker-compose down --volumes --remove-orphans
.PHONY: down

test: 
	go test -v ./...
.PHONY: test