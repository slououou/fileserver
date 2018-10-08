VERSION ?= $(shell ./version.sh)
DOCKER_REPO="fileserver" #change it to your own repo 
DOCKER_IMAGE=$(DOCKER_REPO):$(VERSION)

build:
	@sudo go build -o fileserver

docker-image: 
	@sudo docker build . -t $(DOCKER_IMAGE)


