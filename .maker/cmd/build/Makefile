include $(CURDIR)/../../lib/default.mk

###################
#╔═══════════════╗#
#║ BUILD ACTIONS ║#
#╚═══════════════╝#
###################

include $(GWD)/deploy/token.env

OS := linux darwin
ARCH := amd64 arm64

binaries:
	echo "Building binaries thetiptop in .build/bin"
	rm -rf $(GWD)/.build/bin/*
	for os in $(OS); do \
		for arch in $(ARCH); do \
			ldflags="-s -w -X github.com/kodmain/thetiptop/project/config/config.BUILD_VERSION=$$VERSION -X github.com/kodmain/thetiptop/project/config/config.BUILD_COMMIT=$$(git rev-parse --short HEAD)"; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -buildvcs=false -ldflags="$$ldflags" \
			-o $(GWD)/.build/bin/thetiptop-$$os-$$arch $(GWD)/project/cmd/main.go; \
			chmod +x $(GWD)/.build/bin/thetiptop-$$os-$$arch; \
		done; \
	done

docker: ## Build docker images
	docker build -t thetiptop -f $(GWD)/deploy/Dockerfile $(GWD)/.


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

server: ## Build NomadServer
	echo "Build NomadServer"
	aws --profile=kodmain s3 ls > /dev/null 2>&1 && { echo "Déjà connecté à AWS SSO."; } || { echo "Connexion à AWS SSO requise."; aws sso login --profile=kodmain; }
	terraform -chdir=$(GWD)/deploy/server init -var="github_token=$(GH_TOKEN)"
	terraform -chdir=$(GWD)/deploy/server plan -var="github_token=$(GH_TOKEN)"
	terraform -chdir=$(GWD)/deploy/server apply -auto-approve -var="github_token=$(GH_TOKEN)"

server-destroy: ## Destroy NomadServer
	echo "Destroy NomadServer"
	aws --profile=kodmain s3 ls > /dev/null 2>&1 && { echo "Déjà connecté à AWS SSO."; } || { echo "Connexion à AWS SSO requise."; aws sso login --profile=kodmain; }
	terraform -chdir=$(GWD)/deploy/server apply -auto-approve -destroy -var="github_token=$(GH_TOKEN)"