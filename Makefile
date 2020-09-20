.PHONY: new-version
new-version:
	$(if $(shell which npx),,\
		$(error No npx found in PATH, please install NodeJS))
	$(if $(shell which standard-version),,\
		$(error No standard-version found in PATH, install with: \
			npm install -g standard-version))

	npx standard-version
