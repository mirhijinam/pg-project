up: 
	docker-compose up --build -d && docker-compose logs -f
.PHONY: up

down: 
	docker-compose down --volumes --remove-orphans
.PHONY: down

test: 
	go test -v ./...
.PHONY: test


