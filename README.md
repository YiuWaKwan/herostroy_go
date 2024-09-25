## 线上地址
```
http://cdn0001.afrxvk.cn/hero_story/demo/step040/index.html?serverAddr=127.0.0.1:12345
账号/密码：123
```

```
待解决问题
不对外暴露内部的变量，增加一个 LazySaveRecord，将 SetLastUpdateTime 和 GetLastUpdateTime 移到这里
去除 LazySaveObj 接口上的无用函数，修改 save_or_update.go 代码，保存 LazySaveRecord

switch currLso.(type){}需考虑扩展性
```


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