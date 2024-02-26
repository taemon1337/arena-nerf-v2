IMAGE_REPO ?= taemon1337/arena-game
IMAGE_TAG ?= 0.0.1
DEV_IMAGE ?= arena-game-ui:dev
NET ?= arena-nerf_serf
PORT ?= 8080
HOST ?= $(shell hostname -f)
UID ?= $(shell id -u)
GID ?= $(shell id -g)
PWD ?= $(pwd)
HOME ?= /sonicdev/home

# run dev server (for live coding)
dev:
	docker run --rm -it \
		-p ${PORT}:${PORT} \
		-e NODE_TLS_REJECT_UNAUTHORIZED=0 \
		-e VITE_API=http://127.0.0.1:8000/api/v1 \
		--net ${NET} --user ${UID}:${GID} \
		-v ${HOME}/.kube:/home/node/.kube:ro \
		-v ${PWD}:${PWD} \
		-w ${PWD} "${DEV_IMAGE}" \
		npm run dev

# exec into dev container (for running npm and other tasks)
devshell:
	docker run --rm -it -p ${PORT}:${PORT} --net ${NET} --user ${UID}:${GID} -v ${HOME}/.kube:/home/node/.kube:ro -v ${PWD}:${PWD} -w ${PWD} --entrypoint=/bin/sh "${DEV_IMAGE}"

# build dev image (since we use multistage build we cannot use the deployable image)
build-dev:
	docker build . -f Dockerfile.dev -t "${DEV_IMAGE}"

# build deployable image
build:
	docker build . -t "${IMAGE_REPO}:${IMAGE_TAG}"

# push the built image
push:
	docker push "${IMAGE_REPO}:${IMAGE_TAG}"

# run deployable image locally
server:
	docker run --rm --user ${UID}:${GID} -e VITE_IMAGE_API="http://${HOST}:8090" -v ${HOME}/.kube:/home/.kube:ro -e KUBECONFIG=/home/.kube/config -p ${PORT}:${PORT} ${IMAGE_REPO}:${IMAGE_TAG}
