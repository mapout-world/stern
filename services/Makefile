.PHONY: proto

proto:
	@rm -r proto/**/*.go && \
		$$HOME/.protobuf/bin/protoc proto/**/*.proto \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative