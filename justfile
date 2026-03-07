check-tsp:
	@command -v tsp >/dev/null 2>&1 || { \
	  echo "Error: 'tsp' not found — run 'make install-tsp'"; \
	  exit 1; \
	}

install-tsp:
	npm install -g @typespec/compiler

compile-tsp: check-tsp
	tsp compile ./resources/api/tsp/main.tsp


check-oapi:
	@command -v go tool oapi-codegen --help >/dev/null 2>&1 || { \
	  echo "Error: 'oapi-codegen' not found — run 'make install-oapi'"; \
	  exit 1; \
	}

install-oapi:
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@5b832190e39a030ac999279fbb2d983f2cc73606

generate-oapi: check-oapi
	go tool oapi-codegen -config ./resources/api/oapi_codegen/config.yaml ./resources/api/openapi.yaml
	go mod tidy

tidy:
	@command go mod tidy