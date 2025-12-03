.PHONY: build-rust build-go build-viz dev-up dev-down

build-rust:
	cargo build --workspace

build-go:
	cd go-services/auth && go build -o ../../bin/auth
	cd go-services/audit && go build -o ../../bin/audit

build-viz:
	cd viz && npm run build

dev-up:
	docker compose -f ops/docker-compose.yml up -d

dev-down:
	docker compose -f ops/docker-compose.yml down
