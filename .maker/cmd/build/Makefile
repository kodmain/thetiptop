include $(CURDIR)/../../lib/default.mk

####################
#╔════════════════╗#
#║ SERVER ACTIONS ║#
#╚════════════════╝#
####################

api: ## Build api
	echo "Building binaries thetiptop in .build/bin"
	rm -rf $(GWD)/.build/api/*
	ldflags="-s -w -X github.com/kodmain/thetiptop/api/config/config.BUILD_VERSION=dev -X github.com/kodmain/thetiptop/api/config/config.BUILD_COMMIT=$$(git rev-parse --short HEAD)"; \
	CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -buildvcs=false -ldflags="$$ldflags" -o $(GWD)/.build/api/thetiptop $(GWD)/api/cmd/main.go && chmod +x $(GWD)/.build/api/thetiptop;

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
	swag init -d $(GWD)/api -g cmd/main.go -o $(GWD)/api/internal/docs
	chmod -R 644 $(GWD)/api/internal/docs/*
