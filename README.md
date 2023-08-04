# Context HelloWorld with gRPC Go

项目中大部分内容出自 grpc-go 项目中的 HelloWorld 示例代码。

做了简单的修复：

- 在 client 中插入 metadata
- 在 server 中打印超时时间和 key-val 内容

## 构建

```bash
make
```

## 运行

```bash
./server
./client
```
