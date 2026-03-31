# common

## Prerequisites

Cần cài đặt một số công cụ dưới đây để làm việc với common.

- Golang
- Nodejs (kèm npm, npx)

- Buf cli (v1.57.0)

```bash
go install github.com/bufbuild/buf/cmd/buf@v1.57.0
```

- gRPC plugins

```bash
# protoc-gen-go (v1.36.0)
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0

# protoc-gen-go-grpc (v1.5.1)
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

# protoc-gen-es (v2.7.0), dùng để generate ts/js code dùng cho nodejs hoặc web
npm install -g @bufbuild/protoc-gen-es@v2.7.0
```

Cách sử dụng `buf cli` để generate code cho backend services:

1.  Định nghĩa services trong `api/proto` với cấu trúc được định nghĩa tương tự như `core` đã có sẵn.
    Thư mục `api/vendor` chứa các dependency cơ bản để dùng cho proto, vì không có internet nên vendor về local
    (dùng `buf cli` để lấy deps về và export ra). Khi định nghĩa message có thể validate bằng [protovalidate](https://protovalidate.com/schemas/standard-rules/):
    ````proto
    import "buf/validate/validate.proto";
    ...
        message Example {
            string id = 1 [(buf.validate.field).string.uuid = true];
        }
        ```
    ````
2.  Dùng `buf cli` để generate code vào `/pkg/api/go`, các file `buf.genn.yaml` và `buf.yaml` thể hiện cấu trúc của protos và cấu hình generate code.
    ```bash
    # generate proto
    buf generate
    ```
3.  Import code từ `/pkg/api/<project>` để sử dụng.
