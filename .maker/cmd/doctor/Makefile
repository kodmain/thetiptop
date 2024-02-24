include $(CURDIR)/../../lib/default.mk

###########################
#╔═══════════════════════╗#
#║ DOCUMENTATION ACTIONS ║#
#╚════════════════=══════╝#
###########################

fix-perms: ## Fix permissions
	find $(GWD) -type d -print0 | xargs -0 chmod 0775
	find $(GWD) -type f -print0 | xargs -0 chmod 0664