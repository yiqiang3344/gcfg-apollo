package gcfg_apollo

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/net/context"
)

type AdapterApollo struct {
	Config          config.AppConfig
	Client          agollo.Client
	KeyNamespaceMap map[string]string
}

func (a *AdapterApollo) Available(ctx context.Context, resource ...string) (ok bool) {
	//默认namespace的key的数量大于0则表示可用
	ok = a.Client.GetDefaultConfigCache().EntryCount() > 0
	return
}
func (a *AdapterApollo) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	keys := gstr.Split(pattern, ".")
	fKey := keys[0]
	cfg := a.Client.GetConfig(storage.GetDefaultNamespace())
	if n, ok := a.KeyNamespaceMap[fKey]; ok {
		cfg = a.Client.GetConfig(n)
	}
	j, err := gjson.LoadContentType(gjson.ContentTypeProperties, cfg.GetContent(), true)
	if err != nil {
		return
	}
	//最后一层key要转成小写，兼容gjson处理properties中的奇葩逻辑
	var _keys []string
	if len(keys) > 1 {
		for k, v := range keys {
			if k == len(keys)-1 {
				continue
			}
			_keys = append(_keys, v)
		}
	}
	_keys = append(_keys, gstr.ToLower(keys[len(keys)-1]))
	value = j.Get(gstr.Join(_keys, ".")).Val()
	return
}
func (a *AdapterApollo) Data(ctx context.Context) (data map[string]interface{}, err error) {
	var cfgContentArr []string
	for _, v := range gstr.Split(a.Config.NamespaceName, ",") {
		cfg := a.Client.GetConfig(v)
		cfgContentArr = append(cfgContentArr, cfg.GetContent())
	}
	j, err := gjson.LoadContentType(gjson.ContentTypeProperties, gstr.Join(cfgContentArr, "\n"), true)
	if err != nil {
		return
	}
	data = j.Var().Map()
	return
}

func CreateAdapterApollo(ctx context.Context) (adapter *AdapterApollo, err error) {
	adapter = &AdapterApollo{}
	apolloCfg, err := gcfg.Instance().Get(ctx, "apollo")
	if err != nil {
		return
	}
	apolloCfgMap := apolloCfg.Map()
	adapter.Config = config.AppConfig{
		AppID:         apolloCfgMap["AppID"].(string),
		Cluster:       apolloCfgMap["Cluster"].(string),
		IP:            apolloCfgMap["IP"].(string),
		NamespaceName: apolloCfgMap["NamespaceName"].(string),
		MustStart:     apolloCfgMap["MustStart"].(bool),
	}
	adapter.Client, err = agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return &adapter.Config, nil
	})
	if err != nil {
		return
	}
	//设置key的namespace的map，用于处理非默认namespace的key
	adapter.KeyNamespaceMap = map[string]string{}
	if _map, ok := apolloCfgMap["apolloNotDefaultNamespaceKeyMap"].(map[string]interface{}); ok {
		for k, _v := range _map {
			for _, _v1 := range gstr.Split(gconv.String(_v), ",") {
				adapter.KeyNamespaceMap[_v1] = k
			}
		}
	}
	return
}
