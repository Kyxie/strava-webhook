APP_NAME := strava-webhook
DOCKER_IMAGE := kyxie/$(APP_NAME)
PLATFORM := linux/amd64

# ===== Go Commands =====
run:
	go run main.go

build:
	go build -o $(APP_NAME) .

fmt:
	go fmt ./...

test:
	go test ./...

# ===== Docker Commands =====
docker-build:
	docker build -t $(DOCKER_IMAGE):latest .

docker-cross-build:
	docker buildx build \
		--platform $(PLATFORM) \
		-t $(DOCKER_IMAGE):latest \
		--load .

docker-cross-push:
	docker buildx build \
		--platform $(PLATFORM) \
		-t $(DOCKER_IMAGE):latest \
		--push .