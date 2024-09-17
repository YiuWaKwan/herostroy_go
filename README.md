
## 安装 Google Protobuf Golang 插件

```
go get google.golang.org/protobuf
```

安装完成之后，会在 gopath/bin 目录里多出来一个 protoc-gen-go.exe 文件。
这个 protoc-gen-go.exe 直接是用不到的，但是间接的会被 protoc.exe 调用到……

通过 GameMsgProtocol.proto 生成 go 代码

```
protoc --go_out=. .\GameMsgProtocol.proto
```

