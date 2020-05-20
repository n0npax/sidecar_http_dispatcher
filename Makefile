EXECUTABLES := poetry curl isort bandit
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))


.PHONY: lint
lint:
	sed -i 's/[ \t]*$$//' $(shell find . -name "*.md")
	sed -i 's/[ \t]*$$//' $(shell find . -name "*.py")
	isort $(shell find . -name "*.py") -qy
	bandit -r . -x $(shell find . -name "test_*.py" | tr "\n" ",")
	poetry run black .

cover: test codecov_upload.sh
ifeq ("$(wildcard .codecov_token)","")
	./codecov_upload.sh -c -t "$${CODECOV_TOKEN}" -C "$${COMMIT_SHA}" -Z
else
	./codecov_upload.sh -c -t "$(shell cat .codecov_token)" -Z
endif


.PHONY: test
test:
	poetry run pytest -vs --cov-report=xml --cov . .


.PHONY: build
build:  clean test cover lint
	poetry $@

.PHONY: clean
clean:
	rm -fr dist cli.egg-info poetry.lock codecov_upload.sh || true
	find . -name __pycache__ | xargs rm -fr
	find . -name '*.pyc' -delete

.PHONY: full-clean
full-clean: clean
	rm $${HOME}/.config/lime-comb/ -fr || true
	rm $${HOME}/.local/share/lime-comb/ -fr || true


codecov_upload.sh:
	curl -s https://codecov.io/bash -o codecov_upload.sh
	chmod +x ./codecov_upload.sh
