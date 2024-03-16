include $(shell readlink -f $(CURDIR)/../../lib/default.mk)

##################
#╔══════════════╗#
#║ TEST ACTIONS ║#
#╚══════════════╝#
##################
# Documentation: https://www.atlassian.com/continuous-delivery/software-testing/types-of-software-testing
##################

units: ## Unit tests are very low level and close to the source of an application.
	gotestsum -f github-actions -- -v $$(go list $(GWD)/... | grep -vE "vendor") -coverprofile=$(GWD)/coverage.out -covermode=atomic
integration: ## Integration tests verify that different modules or services used by your application work well together. 
functional: ## Functional tests focus on the business requirements of an application.
end-to-end: ## End-to-end testing replicates a user behavior with the software in a complete application environment.    
acceptance: ## Acceptance tests are formal tests that verify if a system satisfies business requirements. 
performance: ## Performance tests evaluate how a system performs under a particular workload.
smoke: ## Smoke tests are basic tests that check the basic functionality of an application.