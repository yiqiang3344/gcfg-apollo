## goframe的apollo配置适配器
对 `github.com/gogf/gf/contrib/config/apollo/v2` 做了优化，可同时支持多个不同格式的namespace.

### 配置文件示例

```yaml
apollo:
  AppID: "test"
  Cluster: "dev"
  IP: "http://localhost:8080"
  NamespaceName: "application,test.yaml"
  IsBackupConfig: false
  BackupConfigPath: ""
  Secret: ""
  SyncServerTimeout: 0
  MustStart: true
  Watch: true
```

### 依赖
`github.com/apolloconfig/agollo/v4`

### Usage

```bash
go get -u github.com/yiqiang3344/gcfg-apollo@latest
```

在服务启动最开始的地方加上如下代码：

```
adapter, err := gcfg_apollo.CreateAdapterApollo(gctx.New())
if err != nil {
    panic(err)
}
gcfg.Instance().SetAdapter(adapter)
```

具体参考`demo/main.go`，修改`demo/config.yaml`中的apollo配置之后，运行：
```bash
cd demo
GF_GCFG_FILE=config.yaml gf run main.go
```

### 注意
请确保配置中`NamespaceName`的第一个namesapce有发布过，否则会出现`panic: start failed cause no config was read`的错误。