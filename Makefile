test:
	go test ./...
start-all:
	docker-compose down
	docker-compose up &
	sleep 30
	go run cmd/api/http.go & go run cmd/consumer/persist.go && fg