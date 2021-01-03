OS_TARGETS=windows linux
ARCH_TARGETS=386 amd64
OUTPUTS=out/tor-tools-darwin-amd64 $(foreach GOOS,$(OS_TARGETS),$(addprefix out/tor-tools-$(GOOS)-,$(ARCH_TARGETS))) out/tor-tools-linux-arm64

.PHONY: all
all: clean $(OUTPUTS)

define compileTarget =
out/tor-tools-$(1)-$(2): out/
	env GOOS=$(1) GOARCH=$(2) go build -o ./out/tor-tools-$(1)-$(2) ./cmd

endef

$(foreach GOOS,$(OS_TARGETS),$(foreach GOARCH,$(ARCH_TARGETS),$(eval $(call compileTarget,$(GOOS),$(GOARCH)))))

$(eval $(call compileTarget,darwin,amd64))
$(eval $(call compileTarget,linux,arm64))

out/:
	mkdir ./out

.PHONY: test
test:
	@true

.PHONY: clean
clean:
	rm -rf ./out/*
