# gRPC File Upload/Download Service in Go

A high-performance file upload/download server built with gRPC in Go, including:

- Stream-based Upload & Download RPCs
- REST API via gRPC-Gateway
- Cobra-based CLI client with progress indicators
- Benchmarking to compare gRPC vs REST

---

## 🛠️ Requirements

- Go 1.20+
- `protoc` (Protocol Buffers compiler)
- `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-grpc-gateway`
- `make`
- curl (for REST benchmarking)

---

## 📦 Install Dependencies

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

Make sure `GOPATH/bin` is in your `$PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

---

## 📂 Project Structure

```
.
├── proto/                # .proto definitions
├── pb/                   # generated code
├── cmd/                  # CLI client using Cobra
├── server/             # server implementation
├── uploads/              # uploaded/downloaded files
├── third_party/google/api/annotations.proto
├── Makefile
├── main.go
└── README.md
```

---

## 🧰 Make Targets

### ✅ Compile Protobuf

```bash
make proto
```

### ✅ Run gRPC Server (with gRPC-Gateway for REST)

```bash
make run-server
```

### ✅ Run Client Rest

```bash
make run-client-server
```

---

## 🚀 Usage

### Upload a file via CLI

```bash
make run-server
go run client/main.go upload myfile.pdf
```

### Download the same file

```bash
make run-server
go run client/main.go download myfile.pdf downloads/download_
```

You will see progress bars for both upload and download.

---

## 🧪 Benchmark gRPC vs REST

### Step 1: Generate a big test file (e.g., 100MB)

```bash
dd if=/dev/urandom of=bigfile.bin bs=1M count=100
```

### Step 2: Benchmark gRPC Upload (via CLI)

```bash
time go run client/main.go upload bigfile.bin
```

### Step 3: Benchmark REST Upload (via curl)

using base64 if chunks are JSON bytes:

```bash
base64 -i bigfile.bin -o base64.txt
echo '{"filename": "bigfile.bin", "file": "'"$(< base64.txt)"'"}' > upload.json
curl -X POST http://localhost:8080/v1/upload \
  -H "Content-Type: application/json" \
  --data-binary @upload.json
```

### Step 4: Compare execution time and CPU/memory usage

---

## 🔧 Notes

- gRPC is significantly faster and more efficient for large binary data due to protocol buffers.
- REST JSON encoding (especially base64 for bytes) adds overhead.
- This project demonstrates how to bridge both worlds (gRPC + REST) using gRPC-Gateway.

---

## 🧼 Clean Up

```bash
make clean
```

---

## 📜 License

MIT
