SHELL := /bin/bash

.PHONY: run
run:
	../../bin/proaction scan --check outdated-action --quiet --out ./actual.yaml ./workflow.yaml || true
	diff -B expected.yaml actual.yaml