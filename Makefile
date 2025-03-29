build:
	go build -o bin/backend/gateway
	
serve: build
	ENV=development ./bin/backend/gateway

serve-dev-test: build
	ENV=development_test go test -v ./...

docker-push:
	docker build -t antonyinjila/backend/gateway:latest --build-arg ENV=docker .
	docker push antonyinjila/backend/gateway:latest

docker-run:
	docker run -p 8001:8001 ENV=docker antonyinjila/backend/gateway:latest

docker-test:
	ENV=docker_test go test -v ./...