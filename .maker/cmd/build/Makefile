include $(CURDIR)/../../lib/default.mk

####################
#╔════════════════╗#
#║ SERVER ACTIONS ║#
#╚════════════════╝#
####################

OS := linux darwin
ARCH := amd64 arm64

binary: ## Build binary
	echo "Building binaries thetiptop in .build/bin"
	rm -rf $(GWD)/.build/bin/*
	ldflags="-s -w -X github.com/kodmain/thetiptop/project/config/config.BUILD_VERSION=dev -X github.com/kodmain/thetiptop/project/config/config.BUILD_COMMIT=$$(git rev-parse --short HEAD)"; \
	CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -buildvcs=false -ldflags="$$ldflags" -o $(GWD)/.build/bin/thetiptop $(GWD)/project/cmd/main.go && chmod +x $(GWD)/.build/bin/thetiptop;

docker: ## Build docker images
	docker build --build-arg BINARY_VERSION=dev -t ghcr.io/kodmain/thetiptop:dev -f $(GWD)/deploy/Dockerfile $(GWD)/.

certificate: ## Generate ssl certificate
	if [ ! -f $(GWD)/.build/certs/server.key ] || [ ! -f $(GWD)/.build/certs/server.crt ]; then \
		rm -rf $(GWD)/.build/certs; \
		mkdir -p $(GWD)/.build/certs; \
		openssl genrsa -out $(GWD)/.build/certs/server.key 4096; \
		openssl req -new -key $(GWD)/.build/certs/server.key -out $(GWD)/.build/certs/server.csr -sha512 -subj "/CN=localhost"; \
		openssl x509 -req -days 3650 -in $(GWD)/.build/certs/server.csr -signkey $(GWD)/.build/certs/server.key -out $(GWD)/.build/certs/server.crt -sha512; \
	fi

swagger: ## Generate swagger documentation
	echo "Generating swagger documentation"
	swag init -d $(GWD)/project -g cmd/main.go -o $(GWD)/project/internal/docs
	chmod -R 644 $(GWD)/project/internal/docs/*
