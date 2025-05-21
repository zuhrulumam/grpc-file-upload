PROTOC        = protoc
PROTO_SRC     = proto
PROTO_OUT     = pb
PROTO_FILES   = $(wildcard $(PROTO_SRC)/*.proto)

# Tools (assumes they are installed and on your $PATH)
GO_OUT        = --go_out=paths=source_relative:$(PROTO_OUT)
GRPC_OUT      = --go-grpc_out=paths=source_relative:$(PROTO_OUT)

.PHONY: all proto clean

# Generate Go code from .proto
proto:
	@echo "Generating gRPC code from .proto files..."
	$(PROTOC) \
		--proto_path=$(PROTO_SRC) \
		$(GO_OUT) \
		$(GRPC_OUT) \
		$(PROTO_FILES)
	@echo "✅ Done generating proto files in $(PROTO_OUT)/"

# Clean up generated files
clean:
	@echo "Cleaning up generated files..."
	rm -f $(PROTO_OUT)/*.pb.go
	@echo "🧹 Clean complete!"
