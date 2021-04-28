TARGET_APP = noise
DOCKER = docker
IMAGE = timour/$(TARGET_APP)
IMAGE_VERSION = latest
EXPOSED_PORTS = 8080:8080
GOCMD = go

.PHONY: build clean

build: $(TARGET_APP)

$(TARGET_APP): *.go
	CGO_ENABLED=0 $(GOCMD) build -a -installsuffix cgo -o $@ .
	$(DOCKER) build -t $(IMAGE):$(IMAGE_VERSION) .
	$(DOCKER) image prune --force

run:
	$(DOCKER) run -it -e DEBUG=$(DEBUG) -p $(EXPOSED_PORTS) $(IMAGE):$(IMAGE_VERSION)

clean:
	rm -f $(TARGET_APP)