MODULE := github.com/go-park-mail-ru/2025_2_VKarmane

# Все пакеты проекта
ALL_PKGS := $(shell go list ./...)

# Пакеты без моков
PKGS := $(shell echo "$(ALL_PKGS)" | grep -vE '/mocks($|/)')

COVER_MODE := atomic
COVER_RAW := coverage.raw.out
COVER_OUT := coverage.out
COVER_HTML := coverage.html

EXCLUDE_FILES_REGEX := \/mocks\/|\/mock_.*\.go

.PHONY: test-cover cover coverage coverage-raw coverhtml clean

test-cover: cover

# Генерация покрытия без моков
cover:
	@echo "Running tests with coverage (excluding mocks)..."
	GOFLAGS= go test -covermode=$(COVER_MODE) -coverprofile=$(COVER_RAW) $(PKGS)
	@echo "Filtering coverage profile..."
	@head -n1 $(COVER_RAW) > $(COVER_OUT)
	@grep -Ev '$(EXCLUDE_FILES_REGEX)' $(COVER_RAW) | tail -n +2 >> $(COVER_OUT)
	@rm -f $(COVER_RAW)
	@echo "Wrote $(COVER_OUT)"

coverage: cover
	@go tool cover -func=$(COVER_OUT) | grep total:

# Генерация покрытия со всеми пакетами (включая моки)
coverage-raw:
	@echo "Running tests with coverage (WITH mocks)..."
	GOFLAGS= go test -covermode=$(COVER_MODE) -coverprofile=$(COVER_OUT) $(ALL_PKGS)
	@go tool cover -func=$(COVER_OUT) | grep total:

coverhtml: cover
	@go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)
	@echo "Wrote $(COVER_HTML)"
	@which xdg-open >/dev/null 2>&1 && xdg-open $(COVER_HTML) || true
	@which open >/dev/null 2>&1 && open $(COVER_HTML) || true

clean:
	rm -f $(COVER_OUT) $(COVER_HTML) $(COVER_RAW)
