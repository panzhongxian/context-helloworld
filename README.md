# Context HelloWorld with gRPC Go

本项目为 《[【Go】透彻理解 context.Context](https://panzhongxian.cn/cn/2023/08/go-context/)》的示例代码，
项目中大部分内容出自 grpc-go 项目中的 HelloWorld 示例代码。

做了简单的修改：

- 在 greeter\_client 中插入 metadata
- 在 greeter\_server 中打印超时时间和 key-val 内容
- greeter\_client\_with\_otel 在 greeter\_client 的基础上，增加OpenTelemetry Trace
- greeter\_server\_with\_otel 在 greeter\_server 的基础上，增加OpenTelemetry Trace

## 构建

```bash
make
```

## 运行

```bash
./server
./client
```
