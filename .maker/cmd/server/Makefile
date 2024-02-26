include $(CURDIR)/../../lib/default.mk

###################
#╔═══════════════╗#
#║ BUILD ACTIONS ║#
#╚═══════════════╝#
###################

include $(GWD)/deploy/token.env

deploy: ## Build NomadServer
	echo "Build NomadServer"
	aws --profile=kodmain s3 ls > /dev/null 2>&1 && { echo "Déjà connecté à AWS SSO."; } || { echo "Connexion à AWS SSO requise."; aws sso login --profile=kodmain; }
	terraform -chdir=$(GWD)/deploy/server init -var="github_token=$(GH_TOKEN)"
	terraform -chdir=$(GWD)/deploy/server plan -var="github_token=$(GH_TOKEN)"
	terraform -chdir=$(GWD)/deploy/server apply -auto-approve -var="github_token=$(GH_TOKEN)"

destroy: ## Destroy NomadServer
	echo "Destroy NomadServer"
	aws --profile=kodmain s3 ls > /dev/null 2>&1 && { echo "Déjà connecté à AWS SSO."; } || { echo "Connexion à AWS SSO requise."; aws sso login --profile=kodmain; }
	terraform -chdir=$(GWD)/deploy/server apply -auto-approve -destroy -var="github_token=$(GH_TOKEN)"