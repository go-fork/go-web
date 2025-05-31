# Go Fork - Go Package Registry

[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-fork/go-web)](https://goreportcard.com/report/github.com/go-fork/go-web)
[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8.svg)](https://golang.org/doc/go1.23)

## Giới thiệu

Go Fork là một registry cho các Go package chất lượng cao, được thiết kế để giúp các nhà phát triển dễ dàng tìm kiếm và sử dụng các package Go cho dự án của họ. Repository này chứa mã nguồn cho trang web [go.fork.vn](https://go.fork.vn), nơi danh sách tất cả các package có sẵn được hiển thị.

## Tính năng

- Hiển thị danh sách các package Go với mô tả chi tiết
- Hỗ trợ Go module metadata thông qua tham số `go-get=1`
- Tìm kiếm và lọc các package theo tên, mô tả hoặc repository
- Sao chép lệnh `go get` nhanh chóng với một cú nhấp chuột
- Liên kết trực tiếp đến repository và trang tài liệu pkg.go.dev

## Sử dụng

Truy cập [go.fork.vn](https://go.fork.vn) để xem danh sách các package và thêm chúng vào dự án của bạn bằng lệnh `go get`:

```bash
go get go.fork.vn/package-name
```

## Đóng góp

Chúng tôi luôn chào đón sự đóng góp từ cộng đồng! Nếu bạn muốn thêm package của mình vào danh sách:

1. Fork repository này
2. Thêm thông tin package của bạn vào file `data/packages.json`
3. Đảm bảo rằng bạn đã thêm đầy đủ thông tin:
   ```json
   {
     "path": "go.fork.vn/your-package",
     "repository": "github.com/your-username/your-repository",
     "description": "Mô tả ngắn gọn về package của bạn",
     "display": true
   }
   ```
4. Tạo pull request đến repository chính

## Cài đặt và chạy ứng dụng

### Yêu cầu
- Go 1.22 trở lên

### Cách cài đặt

```bash
git clone https://github.com/go-fork/go-web.git
cd go-web
go mod tidy
```

### Chạy server

```bash
go run cmd/main.go
```

Server sẽ chạy mặc định trên cổng 8080. Bạn có thể thay đổi cổng bằng cách đặt biến môi trường `PORT`:

```bash
PORT=8081 go run cmd/main.go
```

## Giấy phép

Dự án này được phân phối dưới giấy phép MIT. Xem file [LICENSE](LICENSE) để biết thêm thông tin.
