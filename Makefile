IMAGE ?= taemon1337/arena-nerf
ARM_IMAGE ?= ws2811-builder:latest
VERSION ?= 2.0.1
PWD ?= $(shell pwd)
APP ?= arena-nerf

build:
	go build .

goarmbuild:
	GOOS=linux GOARCH=arm64 go build -o arena-nerf.arm64

armimage:
	docker buildx build --platform linux/arm64 -f Dockerfile.ws2811 -t ${ARM_IMAGE} .

armbuild:
	docker run --rm -it --platform linux/arm64 -e GOOS=linux -e GOARCH=arm64 -v ${PWD}:/usr/src/${APP} -w /usr/src/${APP} ${ARM_IMAGE} go build -o ${APP}.arm64 -buildvcs=false

run:
	./arena-nerf -enable-controller -enable-game-engine -enable-node -enable-sensor -enable-simulation -enable-connector -name test

docker-armbuild:
	docker build --platform linux/arm64 -f Dockerfile.arm64 -t ${IMAGE}:${VERSION}.arm64 .

docker-build:
	docker build -t ${IMAGE}:${VERSION} .

docker-push:
	docker push ${IMAGE}:${VERSION}

docker-up:
	docker compose up

controller:
	docker run --rm -it \
    -v ./logs:/tmp/logs:rw \
		--net host \
		${IMAGE}:${VERSION} \
		-name control \
		-role ctrl \
		-server \
		-mode domination \
		-start \
		-allow-api-actions \
		-logdir /tmp/logs \
		-gametime 1m \
		-expect 4 \
		-tag role=ctrl \
		-team blue \
		-team red \
		-team green \
		-team yellow

