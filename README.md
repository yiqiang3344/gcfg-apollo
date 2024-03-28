## goframe的apollo配置适配器

### 配置文件示例

```yaml
apollo:
  AppID: "yjqtest"
  Cluster: "dev"
  IP: "http://localhost:8080"
  NamespaceName: "application,test.yaml"
  MustStart: true
  apolloNotDefaultNamespaceKeyMap:
    test.yaml: "server,logger,redis,database"
```

说明

- `apolloNotDefaultNamespaceKeyMap`：表示非默认namespace(也就是application)的，需要配置namespace对应的顶级key的列表。
    - 因为goframe的gcfg的Get()方法没有关于namespace的传值，也不想每次读取配置都要把所有namespace的配置都遍历一遍，所以用了这种方法。

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