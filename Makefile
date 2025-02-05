
SHELL=/bin/bash


.PHONY: install
install:
	go install go.osspkg.com/goppy/v2/cmd/goppy@latest
	goppy setup-lib

.PHONY: lint
lint:
	goppy lint

.PHONY: license
license:
	goppy license

.PHONY: build
build:
	goppy build --arch=amd64

.PHONY: tests
tests:
	goppy test

.PHONY: pre-commit
pre-commit: install license lint tests build

.PHONY: ci
ci: pre-commit

run_back:
	go run cmd/fdns/main.go --config=config/config.dev.yaml

nslookup:
	nslookup -port=8053 google.com 127.0.0.1
#	nslookup -port=8053 adstop.org 127.0.0.1
	nslookup -port=8053 yandex.ru 127.0.0.1
	nslookup -port=8053 vk.com 127.0.0.1
#	nslookup -port=8053 dewep.pro 127.0.0.1
#	nslookup -port=8053 dewep.online 127.0.0.1
#	nslookup -port=8053 googleads.github.io 127.0.0.1
#	nslookup -port=8053 logo-net.co.uk 127.0.0.1
	nslookup -port=8053 1-2-3-4.local 127.0.0.1