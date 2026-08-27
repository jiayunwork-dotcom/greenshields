# greenshields: Go HTTP 后端核算 Greenshields 交通流基本图与冲击波波速

用户给出自由流速度 vf 和堵塞密度 kj，核算速度–密度直线、流量–密度抛物线、通行能力点，以及两交通状态间的运动学波速。必须同时成立：v = vf(1−k/kj)，q = k·v 在 k=kj/2 取最大 qm=vf·kj/4，波速 w=(q2−q1)/(k2−k1)。交叉规则：k=0 时 v=vf、q=0；k=kj 时 v=0、q=0；同一 q<qm 的畅通根与拥堵根之和等于 kj；拥堵上溢时下游近堵塞则 w 为负、扰动向上游传播。

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
