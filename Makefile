lint:
	@sh -c "'$(CURDIR)/scripts/lint.sh'"

coverage:
	@sh -c "'$(CURDIR)/scripts/coverage.sh'"

build:
	@sh -c "'$(CURDIR)/scripts/build.sh'"

pipeline:
	make lint && make coverage && make build

.PHONY: lint, coverage, build, pipeline

.DEFAULT_GOAL := pipeline
