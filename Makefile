

up:
	docker-compose up -d --build
down:
	docker-compose down -v
debug:
	docker-compose exec -it go-service /bin/bash
click:
	docker-compose exec -it clickhouse-server clickhouse-client


GOBIN := /Users/mikhail.golyachkov/go/bin
GO_OUT_PATH := ./go-service/proto

.create_go_out:
	mkdir -p $(GO_OUT_PATH)
proto-go: .create_go_out
	protoc \
		--proto_path=./proto \
		--plugin=protoc-gen-go=$(GOBIN)/protoc-gen-go \
		--plugin=protoc-gen-go-grpc=$(GOBIN)/protoc-gen-go-grpc \
		--go_out=paths=source_relative:${GO_OUT_PATH} \
		--go-grpc_out=paths=source_relative:${GO_OUT_PATH} \
		$$(find proto -name '*.proto' | xargs)


PYTHON := ./py-service/.venv/bin/python
PY_OUT := ./py-service

proto-python:
	mkdir -p $(PY_OUT)
	$(PYTHON) -m grpc_tools.protoc \
		-I ./proto \
		--python_out=$(PY_OUT) \
		--pyi_out=$(PY_OUT) \
		--grpc_python_out=$(PY_OUT) \
		$$(find proto -name '*.proto' | xargs)