NAME ?= shadowsocks-ws
BINDIR=bin
GOBUILD=CGO_ENABLED=0 go build -ldflags '-w -s'
# The -w and -s flags reduce binary sizes by excluding unnecessary symbols and debug info

.PHONY: all
all: arm linux macos win64

# cubietrunk plus, CPU: ARMCortex A7, arch: ARMv7-A
# raspberry pi 3B+, CPU: ARM Cortex-A53,  arch: ARMv8-A
.PHONY: arm
arm:
	GOARCH=arm GOARM=7 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

.PHONY: linux
linux:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

.PHONY: macos
macos:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

.PHONY: win64
win64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

.PHONY: releases
releases: clean linux macos win64
	chmod +x $(BINDIR)/$(NAME)-*
	gzip $(BINDIR)/$(NAME)-linux
	gzip $(BINDIR)/$(NAME)-macos
	zip -m -j $(BINDIR)/$(NAME)-win64.zip $(BINDIR)/$(NAME)-win64.exe

.PHONY: clean
clean:
	rm $(BINDIR)/*