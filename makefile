.PHONY: build

apid local start:
	node_modules/.bin/aglio --theme-variables slate -i ./docs/index.apib -s -p 80 -h '0.0.0.0'

