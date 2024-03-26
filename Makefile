include $(shell readlink -f $(CURDIR)/.maker/lib/default.mk)

doc: documentation ## Display all documentation actions

documentation:
	make -C $(CURDIR)/.maker/cmd/doc $(ARGS)

build: ## Display all builds actions
	make -C $(CURDIR)/.maker/cmd/build $(ARGS)

test: ## Display all tests actions
	make -C $(CURDIR)/.maker/cmd/test $(ARGS)

server: ## Display all servers actions
	make -C $(CURDIR)/.maker/cmd/server $(ARGS)

doctor: ## Display all doctor actions
	make -C $(CURDIR)/.maker/cmd/doctor $(ARGS)

update: ## Install/Update vendor
	echo "Update all dependencies"
	go get -u $(GWD)/...
	go mod vendor
	cd $(GWD)/app && flutter pub outdated && flutter pub get

run: ## RUN api without build
	make build certificate > /dev/null
	make build swagger 	   > /dev/null
	go run $(GWD)/api/cmd/main.go -t $(GWD)/.build/certs