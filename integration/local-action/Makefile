SHELL := /bin/bash

.PHONY: run
run:
	../../bin/proaction scan --quiet --out ./actual.yaml ./workflow.yaml || true
	diff -B expected.yaml actual.yaml
