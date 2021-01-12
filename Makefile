ARCH_TARGETS=386 amd64
OUTPUTS=out/tor-tools-darwin-amd64 $(addprefix out/tor-tools-windows-$(GOARCH),$(addsuffix .exe,$(ARCH_TARGETS))) $(addprefix out/tor-tools-linux-,$(ARCH_TARGETS) arm64)

.PHONY: all
all: clean $(OUTPUTS)

define compileTarget =
ifeq ($(1),windows)
out/tor-tools-$(1)-$(2).exe: out/
	env GOOS=$(1) GOARCH=$(2) go build -o ./out/tor-tools-$(1)-$(2).exe ./cmd
else
out/tor-tools-$(1)-$(2): out/
	env GOOS=$(1) GOARCH=$(2) go build -o ./out/tor-tools-$(1)-$(2) ./cmd
endif

endef

$(foreach GOARCH,$(ARCH_TARGETS) arm64,$(eval $(call compileTarget,linux,$(GOARCH))))
$(foreach GOARCH,$(ARCH_TARGETS),$(eval $(call compileTarget,windows,$(GOARCH))))
$(eval $(call compileTarget,darwin,amd64))

out/:
	mkdir ./out

.PHONY: test
test:
	@true

.PHONY: clean
clean:
	rm -rf ./out
