PREFIX ?= /usr/local
BIN    := get-ip
PKG    := ./
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

ZSHRC      := $(HOME)/.zshrc
ALIAS_START := \# >>> get-ip alias >>>
ALIAS_END   := \# <<< get-ip alias <<<

.PHONY: help build install uninstall alias unalias test fmt clean run release formula-sha

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build binary into bin/
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/$(BIN) $(PKG)

install: build ## Install to $(PREFIX)/bin
	install -m 0755 bin/$(BIN) $(PREFIX)/bin/$(BIN)
	@echo "installed to $(PREFIX)/bin/$(BIN)"

uninstall: ## Remove installed binary and alias block
	rm -f $(PREFIX)/bin/$(BIN)
	@if [ -f $(ZSHRC) ] && grep -q '>>> get-ip alias >>>' $(ZSHRC); then \
		sed -i.bak '/>>> get-ip alias >>>/,/<<< get-ip alias <<</d' $(ZSHRC); \
		echo "alias block removed from $(ZSHRC) (backup: $(ZSHRC).bak)"; \
	fi

alias: ## Append guarded `ip` alias to ~/.zshrc
	@if grep -q '>>> get-ip alias >>>' $(ZSHRC) 2>/dev/null; then \
		echo "alias block already present in $(ZSHRC)"; \
	else \
		printf '\n# >>> get-ip alias >>>\nalias ip='\''get-ip'\''\n# <<< get-ip alias <<<\n' >> $(ZSHRC); \
		echo "added alias to $(ZSHRC). run: source ~/.zshrc"; \
	fi

unalias: ## Remove guarded alias block from ~/.zshrc
	@if [ -f $(ZSHRC) ] && grep -q '>>> get-ip alias >>>' $(ZSHRC); then \
		sed -i.bak '/>>> get-ip alias >>>/,/<<< get-ip alias <<</d' $(ZSHRC); \
		echo "alias block removed (backup: $(ZSHRC).bak). run: source ~/.zshrc"; \
	else \
		echo "no alias block found"; \
	fi

test: ## Run tests
	go test ./...

fmt: ## Format and vet
	gofmt -w .
	go vet ./...

clean: ## Remove bin/
	rm -rf bin/

run: ## go run
	go run $(PKG)

release: ## Tag a release and print formula update snippet (usage: make release V=v0.1.0)
	@if [ -z "$(V)" ]; then echo "usage: make release V=vX.Y.Z" >&2; exit 2; fi
	./scripts/release.sh $(V)

formula-sha: ## Print SHA256 of an existing release tarball (usage: make formula-sha V=v0.1.0)
	@if [ -z "$(V)" ]; then echo "usage: make formula-sha V=vX.Y.Z" >&2; exit 2; fi
	@curl -fsSL "https://github.com/vianch/get-ip/archive/refs/tags/$(V).tar.gz" | shasum -a 256 | awk '{print $$1}'
