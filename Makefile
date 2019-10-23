# Super simple makefile

# Base layer dependencies
GO_VERSION=1.13
ALPINE_VERSION=3.10

# TODO wire up versioning consistently...
#VERSION=
TAG=latest # TODO version based
ORG?=docker.io/dhiltgen
DAEMON_IMAGE=$(ORG)/site:$(TAG)
DOCKER?=docker

# Using buildx for multi-arch
#DOCKER_BUILD=$(DOCKER) buildx build --push \
#    --build-arg GO_VERSION=$(GO_VERSION) \
#    --build-arg ALPINE_VERSION=$(ALPINE_VERSION) \
#    --platform linux/amd64,linux/arm/v7
# Using regular build without multi-arch
DOCKER_BUILD=$(DOCKER) build --build-arg GO_VERSION=$(GO_VERSION) --build-arg ALPINE_VERSION=$(ALPINE_VERSION)

DOCKER_BUILDKIT=1
export DOCKER_BUILDKIT


build: daemon

daemon:
	$(DOCKER_BUILD) \
	    -t $(DAEMON_IMAGE) \
	    --target daemon \
	    .


# Doesn't benefit from multi-arch at present
test:
	$(DOCKER) build \
	    --target test \
	    .

# Doesn't benefit from multi-arch at present
coverage:
	$(DOCKER) build \
	    --target coverage \
	    --output type=local,dest=./ \
	    .
	xdg-open ./cover.html

# Only used for non-multi-arch builds
push:
	$(DOCKER) push $(DAEMON_IMAGE)

fmt:
	clang-format -i $$(find . -name \*.proto | grep -v "vendor/")

.PHONY: build daemon test deploy coverage push fmt
