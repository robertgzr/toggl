COMMAND := bin/toggl
GOBUILD_FLAGS := -ldflags "-X github.com/robertgzr/toggl/app.VERSION=$(shell git describe --tags --always) -X github.com/robertgzr/toggl/app.BUILDTIME=$(shell date -u +%Y%m%d.%H%M%S)"

all: $(COMMAND)

clean:
	$(RM) $(COMMAND)
	# git clean -n

build: clean $(COMMAND)
$(COMMAND):
	go build $(GOBUILD_FLAGS) -o $@ 

.PHONY: build
