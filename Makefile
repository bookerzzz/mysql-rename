build:
	mkdir -p "dist/bin"
	GOBIN="$(CURDIR)/dist/bin" go install .

