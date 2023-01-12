.PHONY: *
IMAGE_NAME=docker-url-collector
URL_COLLECTOR_PORT=8091

build:
	DOCKER_BUILDKIT=1 docker build \
	--build-arg http_proxy \
	--build-arg https_proxy \
	--build-arg ftp_proxy \
	--build-arg no_proxy \
	--build-arg=URL_COLLECTOR_PORT=${URL_COLLECTOR_PORT} \
	--tag ${IMAGE_NAME}:latest .

run:
	docker run ${IMAGE_NAME}
