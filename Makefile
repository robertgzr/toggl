
OUTDIR := bin
VERSION := $(shell git describe --tags --always)
GOFLAGS ?=
LDFLAGS := \
	-X main.version=$(VERSION) \
	-X main.date=$(shell date -u +%Y%m%d.%H%M%S)

GOSRC!=find . -name '*.go'
GOSRC+=go.mod go.sum


all: $(OUTDIR)/toggl

build: clean $(OUTDIR)/toggl

$(OUTDIR)/toggl: $(GOSRC)
	go build $(GOFLAGS) \
		-ldflags "$(LDFLAGS)" \
		-o $@

clean:
	$(RM) -r $(OUTDIR)
	$(RM) -r ./dist

PREFIX?=/usr/local
_INSTDIR=$(DESTDIR)$(PREFIX)
BINDIR?=$(_INSTDIR)/bin

install: all
	mkdir -p $(BINDIR)
	install -m755 $(OUTDIR)/toggl $(BINDIR)/toggl
	install -m755 toggl_dmenu.sh $(BINDIR)/toggl_dmenu

uninstall:
	$(RM) $(BINDIR)/toggl
	$(RM) $(BINDIR)/toggl_dmenu


.PHONY: clean install uninstall
