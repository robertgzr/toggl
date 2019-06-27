COMMAND := bin/toggl
GOBUILD_FLAGS := -ldflags "-X main.version=$(shell git describe --tags --always) -X main.date=$(shell date -u +%Y%m%d.%H%M%S)"

all: $(COMMAND)

build: clean $(COMMAND)
$(COMMAND):
	go build $(GOBUILD_FLAGS) -o $@ 

clean:
	$(RM) $(COMMAND)
	$(RM) -r ./dist
	# git clean -n

.PHONY: build clean
