BUILD_VERSION?=0.0.1
BUILD_TIME?=$(shell date)

kind-cluster:
	kind create cluster --name=kxc --config=kind-config.yaml

k8s-deploy:
	kubectl apply -f k8s/ingress/

run-kxc-api:
	(cd kxc-api && go run cmd/main.go)

run-kxc-api-bootstrap:
	(cd kxc-api && go run cmd/main.go bootstrap)

run-kxc-operator:
	(cd kxc-operator && make run ENABLE_WEBHOOKS=false)

build-kxc-cli:
	(cd kxc-cli && go build -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.BuildVersion=$(BUILD_VERSION)'"  -o bin/kxc ./cmd/)

test:
	(cd kxc-operator && make test)
	(cd kxc-api && go test ./...)
	(cd kxc-cli && go test ./...)

mockery:
	(cd kxc-api && mockery --dir "./services" --all --output "./testsupport/mocks")