# ===========================================
# Project config
# ===========================================
APP := gitwo
PKG := ./cmd/gitwo
LDFLAGS := -s -w -X main.version=dev -X main.commit=$$(git rev-parse --short HEAD 2>/dev/null || echo none) -X main.date=$$(date -u +'%Y-%m-%dT%H:%M:%SZ')

# Local build cache: rebuild only if sources changed
SRC_CHECKSUM := $(shell find . -name "*.go" -type f -exec sha256sum {} \; | sort | sha256sum | cut -d' ' -f1)
CHECKSUM_FILE := .build-checksum

# act (run GitHub Actions locally)
ACT_IMAGE    ?= catthehacker/ubuntu:act-22.04
ARCH_FLAG    ?=
EVENT_DRY    ?= .github/events/dispatch.json
EVENT_TAG    ?= .github/events/push_tag.json
SECRETS_FILE ?= .secrets.env

# -------------------------------------------
# Phonies
# -------------------------------------------
.PHONY: build force-build run test clean snapshot snapshot-fast release check flow \
        verify guard-clean-tree guard-tidy-clean guard-generated-clean guard-tests guard-tag-sane \
        prepare-archives act-setup act-dry act-release help

# ===========================================
# Core dev targets
# ===========================================
build:
	@echo "🔍 Checking if rebuild is necessary..."
	@if [ -f "$(CHECKSUM_FILE)" ] && [ -f "bin/$(APP)" ] && [ "$$(cat $(CHECKSUM_FILE))" = "$(SRC_CHECKSUM)" ]; then \
		echo "✅ Build is up to date - no changes detected"; \
		echo "   Binary: bin/$(APP)"; \
		echo "   Checksum: $(SRC_CHECKSUM)"; \
	else \
		echo "🔨 Building $(APP)..."; \
		go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP) $(PKG); \
		echo "$(SRC_CHECKSUM)" > $(CHECKSUM_FILE); \
		echo "✅ Build complete - checksum saved"; \
	fi

force-build:
	@echo "🔨 Force building $(APP)..."
	@go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP) $(PKG)
	@echo "$(SRC_CHECKSUM)" > $(CHECKSUM_FILE)
	@echo "✅ Force build complete - checksum saved"

run: build
	./bin/$(APP) --version

test:
	go test ./...

clean:
	rm -rf bin dist $(CHECKSUM_FILE)

# ===========================================
# Release safety & verification
# ===========================================
# 1) Fail if working tree has uncommitted changes (tracked files)
guard-clean-tree:
	@git diff --quiet || (echo "❌ Working tree is dirty (tracked changes). Commit or stash."; exit 1)
	@git diff --cached --quiet || (echo "❌ You have staged but uncommitted changes."; exit 1)

# 2) go mod tidy must be a no-op (prevents module drift on CI)
guard-tidy-clean:
	@cp go.mod go.mod.bak 2>/dev/null || true
	@cp go.sum go.sum.bak 2>/dev/null || true
	@go mod tidy
	@git diff --quiet -- go.mod go.sum || (echo "❌ 'go mod tidy' changed files. Commit those changes."; git --no-pager diff -- go.mod go.sum; mv -f go.mod.bak go.mod 2>/dev/null || true; mv -f go.sum.bak go.sum 2>/dev/null || true; exit 1)
	@rm -f go.mod.bak go.sum.bak 2>/dev/null || true

# 3) If you have code generators, run them and ensure no diffs.
guard-generated-clean:
	@echo "ℹ️ No generators configured. If you add any, wire them here and check for diffs."
	@true

# 4) Tests
guard-tests:
	@go test ./... >/dev/null || (echo "❌ Tests failed."; exit 1)

# 5) Tag sanity: only allow release when on an exact vX.Y.Z tag that points to HEAD
guard-tag-sane:
	@tag=$$(git describe --tags --exact-match 2>/dev/null || true); \
	if [ -z "$$tag" ]; then echo "❌ Not on an exact tag. Create and push tag vX.Y.Z."; exit 1; fi; \
	case "$$tag" in v[0-9]*.[0-9]*.[0-9]*) ;; *) echo "❌ Tag '$$tag' is not semver (vX.Y.Z)."; exit 1;; esac; \
	echo "✅ On release tag $$tag"

# Aggregate verification (what CI should run pre-release)
verify: guard-clean-tree guard-tidy-clean guard-generated-clean guard-tests
	@echo "✅ verify passed (clean tree, tidy clean, generators clean, tests ok)"

check:
	goreleaser check

# Ensure required files for archives/formula exist
prepare-archives:
	@[ -f LICENSE ]   || (echo "❌ LICENSE missing at repo root"; exit 1)
	@[ -f README.md ] || (echo "❌ README.md missing at repo root"; exit 1)
	@for f in completions/gitwo.bash completions/gitwo.zsh completions/gitwo.fish; do \
	  [ -f $$f ] || { echo "❌ $$f missing"; exit 1; }; \
	done

# ===========================================
# GoReleaser
# ===========================================
# Local build of artifacts (no publish)
snapshot:
	goreleaser release --snapshot --skip=publish --clean

# Faster local dry-run (skip docker/sign/sbom)
snapshot-fast:
	goreleaser release --snapshot --skip=publish --skip=docker --skip=sign --skip=sbom --clean

# Guarded local release (publishes) — CI will run the same steps
release: verify guard-tag-sane prepare-archives
	@if [ -z "$$GITHUB_TOKEN" ]; then echo "ERROR: GITHUB_TOKEN is required for GitHub release"; exit 1; fi
	@if [ -z "$$HOMEBREW_TAP_GITHUB_TOKEN" ]; then echo "ERROR: HOMEBREW_TAP_GITHUB_TOKEN is required for Homebrew tap push"; exit 1; fi
	goreleaser release --clean

# ===========================================
# GitHub Actions via 'act'
# ===========================================
act-setup:
	@mkdir -p .github/events
	@[ -f $(EVENT_DRY) ] || printf '{ "ref": "refs/heads/main", "repository": { "full_name": "gitwohq/gitwo", "name": "gitwo", "owner": { "login": "gitwohq" } }, "inputs": {} }\n' > $(EVENT_DRY)
	@[ -f $(EVENT_TAG) ] || printf '{ "ref": "refs/tags/v0.1.0", "repository": { "full_name": "gitwohq/gitwo", "name": "gitwo", "owner": { "login": "gitwohq" } } }\n' > $(EVENT_TAG)
	@if ! grep -q '^\.secrets\.env$$' .gitignore 2>/dev/null; then echo ".secrets.env" >> .gitignore; fi
	@echo "✅ act setup ready. Put tokens (if needed) into $(SECRETS_FILE)."

act-dry: act-setup
	@echo "▶ Running DRY-RUN workflow locally (no publish)"
	act -W .github/workflows/release.dry.yml \
	    -j goreleaser \
	    -e $(EVENT_DRY) \
	    -P ubuntu-22.04=$(ACT_IMAGE) \
	    $(ARCH_FLAG)

act-release: act-setup
	@echo "▶ Running REAL release workflow locally (publishes if tokens are set)"
	@if [ ! -f "$(SECRETS_FILE)" ]; then echo "Hint: create $(SECRETS_FILE) with HOMEBREW_TAP_GITHUB_TOKEN=... and GITHUB_TOKEN=..."; fi
	act -W .github/workflows/release.yml \
	    -j goreleaser \
	    -e $(EVENT_TAG) \
	    -P ubuntu-22.04=$(ACT_IMAGE) \
	    --secret-file $(SECRETS_FILE) \
	    $(ARCH_FLAG)

# ===========================================
# Flow & Help
# ===========================================
flow:
	@echo '== Local dev flow =='; \
	echo '  1) make check'; \
	echo '  2) make snapshot-fast   # quick dry-run (no docker/sign/sbom)'; \
	echo '  3) make snapshot        # full dry-run (no publish)'; \
	echo; \
	echo '== CI real release =='; \
	echo '  - git tag -a vX.Y.Z -m "vX.Y.Z" && git push origin vX.Y.Z'; \
	echo '  - .github/workflows/release.yml runs GoReleaser and publishes'; \
	echo; \
	echo '== CI simulation (local) =='; \
	echo '  - make act-dry          # runs release.dry.yml via act'; \
	echo '  - make act-release      # runs release.yml via act (needs tokens)'; \
	echo; \
	echo 'See docs: docs/maintainer-playbook.mdx (Releasing).'

help:
	@echo "Make targets:"
	@echo "  build           - build binary if sources changed (cached)"
	@echo "  force-build     - force rebuild binary"
	@echo "  run             - run './bin/$(APP) --version'"
	@echo "  test            - go test ./..."
	@echo "  clean           - remove bin/dist/checksum"
	@echo "  check           - goreleaser config sanity check"
	@echo "  verify          - pre-release guards (clean tree, tidy clean, tests)"
	@echo "  prepare-archives - assert LICENSE/README/completions exist"
	@echo "  snapshot        - local goreleaser snapshot (no publish)"
	@echo "  snapshot-fast   - snapshot without docker/sign/sbom"
	@echo "  release         - local goreleaser release (publishes) [requires tokens] + guards"
	@echo "  act-setup       - create act event payloads + ignore .secrets.env"
	@echo "  act-dry         - run DRY GitHub workflow locally (no publish)"
	@echo "  act-release     - run REAL GitHub workflow locally (publishes if tokens exist)"
	@echo "  flow            - print the recommended order of commands"
