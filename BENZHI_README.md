# greenshields：Greenshields 交通流基本图（速度-密度-流量与冲击波波速）

Go 同进程托管静态前端（web/）与 JSON API；启动后访问 http://localhost:8080/。

## 构建 / 运行 / 测试

```text
go build ./...
go run . -http :8080
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
