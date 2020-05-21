SHELL := /bin/bash

.PHONY: run
run:
	../../bin/proaction scan --check recommendations --quiet --out ./actual.yaml ./workflow.yaml || true
	diff -B expected.yaml actual.yaml
