include $(shell readlink -f $(CURDIR)/.maker/lib/default.mk)

doc: documentation ## Display all documentation actions

documentation:
	make -C $(CURDIR)/.maker/cmd/doc $(ARGS)

build: ## Display all builds actions
	make -C $(CURDIR)/.maker/cmd/build $(ARGS)

test: ## Display all tests actions
	make -C $(CURDIR)/.maker/cmd/test $(ARGS)

doctor: ## Display all doctor actions
	make -C $(CURDIR)/.maker/cmd/doctor $(ARGS)

update: ## Install/Update vendor
	echo "Update all dependencies"
	go get -u $(GWD)/...
	go mod vendor

run: ## Run the project localy
	make build certificate > /dev/null
	make build swagger > /dev/null
	go run project/cmd/main.go -t $(GWD)/.build/certs
