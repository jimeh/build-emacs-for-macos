.PHONY: new-version
new-version: check-npx
	npx standard-version

.PHONY: next-version
next-version: check-npx
	npx standard-version --dry-run

.PHONY: check-npx
check-npx:
	$(if $(shell which npx),,\
		$(error No npx execuable found in PATH, please install NodeJS))
