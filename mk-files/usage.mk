define MIGRATEDOWN
BEGIN;

DELETE FROM whitelist WHERE version = '$(BUMPED_VERSION)';

COMMIT;
endef

export MIGRATEDOWN
update-whitelist:
	$(eval CC_CLI_SERVICE=$(shell mktemp -d)/cc-cli-service)
	
	git clone git@github.com:confluentinc/cc-cli-service.git $(CC_CLI_SERVICE) && \
	cd $(CC_CLI_SERVICE) && \
	git checkout -b update-whitelist-$(BUMPED_VERSION) && \
	make db-migrate-create NAME=$(BUMPED_VERSION)
	
	go run -ldflags "-X main.version=$(BUMPED_VERSION)" cmd/usage/main.go > $$(find $(CC_CLI_SERVICE)/db/migrations/ -name "*_$(BUMPED_VERSION).up.sql")
	echo "$$MIGRATEDOWN" > $$(find $(CC_CLI_SERVICE)/db/migrations/ -name "*_$(BUMPED_VERSION).down.sql")
	
	cd $(CC_CLI_SERVICE) && \
	make db-migrate-up && \
	git add . && \
	git commit -m "update whitelist for $(BUMPED_VERSION)" && \
	git push origin update-whitelist-$(BUMPED_VERSION) && \
	gh pr create -B master --title "Update whitelist for $(BUMPED_VERSION)" --body ""

promote:
	$(eval DIR=$(shell mktemp -d))
	$(eval CC_CLI_SERVICE=$(DIR)/cc-cli-service)
	
	git clone git@github.com:confluentinc/cc-cli-service.git $(CC_CLI_SERVICE) && \
	cd $(CC_CLI_SERVICE) && \
	git checkout -b promote-$(BUMPED_VERSION) && \
	export VAULT_ADDR=https://vault.cireops.gcp.internal.confluent.cloud; vault login -method=oidc -path=okta/ && \
	halctl --context prod --vault-oidc-role halyard-prod --vault-token "$$(cat ~/.vault-token)" --vault-login-path "auth/app/prod/login" release service environment version list cc-cli-service stag | sed -n 's/\*.*\(0.[0-9]*.0\).*/\1/p' > $(DIR)/version.txt && \
	for env in devel prod; do \
		halctl --context prod --vault-oidc-role halyard-prod --vault-token "$$(cat ~/.vault-token)" --vault-login-path "auth/app/prod/login" release service environment version list cc-cli-service $${env} | grep $$(cat $(DIR)/version.txt) | grep -o -E '[0-9]+' | head -1 > $(DIR)/$${env}.txt; \
		sed -i '' "s/installedVersion: \"[0-9]*\"/installedVersion: \"$$(cat $(DIR)/$${env}.txt)\"/" .deployed-versions/$${env}.yaml; \
	done && \
	git commit -am "promote devel to $$(cat $(DIR)/devel.txt) and promote prod to $$(cat $(DIR)/prod.txt) for $(BUMPED_VERSION)" && \
	git push origin promote-$(BUMPED_VERSION) && \
	gh pr create -B master --title "Promote devel and prod for $(BUMPED_VERSION)" --body ""
