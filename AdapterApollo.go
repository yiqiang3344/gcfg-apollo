package gcfg_apollo

import (
	"context"
	"github.com/apolloconfig/agollo/v4"
	apolloConfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type Config struct {
	AppID             string `v:"required"` // See apolloConfig.Config.
	IP                string `v:"required"` // See apolloConfig.Config.
	Cluster           string `v:"required"` // See apolloConfig.Config.
	NamespaceName     string // See apolloConfig.Config.
	IsBackupConfig    bool   // See apolloConfig.Config.
	BackupConfigPath  string // See apolloConfig.Config.
	Secret            string // See apolloConfig.Config.
	SyncServerTimeout int    // See apolloConfig.Config.
	MustStart         bool   // See apolloConfig.Config.
	Watch             bool   // Watch watches remote configuration updates, which updates local configuration in memory immediately when remote configuration changes.
}

type AdapterApollo struct {
	config Config
	client agollo.Client
	value  *g.Var // Configmap content cached. It is `*gjson.Json` value internally.
}

func CreateAdapterApollo(ctx context.Context) (adapter *AdapterApollo, err error) {
	apolloCfg, err := gcfg.Instance().Get(ctx, "apollo")
	if err != nil {
		return
	}
	apolloCfgMap := apolloCfg.Map()
	config := Config{
		AppID:             apolloCfgMap["AppID"].(string),
		IP:                apolloCfgMap["IP"].(string),
		Cluster:           apolloCfgMap["Cluster"].(string),
		NamespaceName:     apolloCfgMap["NamespaceName"].(string),
		IsBackupConfig:    apolloCfgMap["IsBackupConfig"].(bool),
		BackupConfigPath:  apolloCfgMap["BackupConfigPath"].(string),
		Secret:            apolloCfgMap["Secret"].(string),
		SyncServerTimeout: gconv.Int(apolloCfgMap["SyncServerTimeout"]),
		MustStart:         apolloCfgMap["MustStart"].(bool),
		Watch:             apolloCfgMap["Watch"].(bool),
	}
	// Data validation.
	err = g.Validator().Data(config).Run(ctx)
	if err != nil {
		return nil, err
	}
	if config.NamespaceName == "" {
		config.NamespaceName = storage.GetDefaultNamespace()
	}
	adapter = &AdapterApollo{
		config: config,
		value:  g.NewVar(nil, true),
	}
	adapter.client, err = agollo.StartWithConfig(func() (*apolloConfig.AppConfig, error) {
		return &apolloConfig.AppConfig{
			AppID:             config.AppID,
			Cluster:           config.Cluster,
			NamespaceName:     config.NamespaceName,
			IP:                config.IP,
			IsBackupConfig:    config.IsBackupConfig,
			BackupConfigPath:  config.BackupConfigPath,
			Secret:            config.Secret,
			SyncServerTimeout: config.SyncServerTimeout,
			MustStart:         config.MustStart,
		}, nil
	})
	if err != nil {
		return
	}
	if config.Watch {
		adapter.client.AddChangeListener(adapter)
	}
	return
}

func (a *AdapterApollo) Available(ctx context.Context, resource ...string) (ok bool) {
	if len(resource) == 0 && !a.value.IsNil() {
		return true
	}
	var namespace = gstr.Split(a.config.NamespaceName, ",")[0]
	if len(resource) > 0 {
		namespace = resource[0]
	}
	return a.client.GetConfig(namespace) != nil
}

func (a *AdapterApollo) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	if a.value.IsNil() {
		if err = a.updateLocalValue(ctx); err != nil {
			return nil, err
		}
	}
	return a.value.Val().(*gjson.Json).Get(dealPattern(pattern)).Val(), nil
}

func (a *AdapterApollo) Data(ctx context.Context) (data map[string]interface{}, err error) {
	if a.value.IsNil() {
		if err = a.updateLocalValue(ctx); err != nil {
			return nil, err
		}
	}
	return a.value.Val().(*gjson.Json).Map(), nil
}

func (c *AdapterApollo) OnChange(event *storage.ChangeEvent) {
	_ = c.updateLocalValue(gctx.New())
}

// OnNewestChange is called when any config changes.
func (c *AdapterApollo) OnNewestChange(event *storage.FullChangeEvent) {
	// Nothing to do.
}

func (c *AdapterApollo) updateLocalValue(ctx context.Context) (err error) {
	var j = gjson.New(nil)
	for _, v := range gstr.Split(c.config.NamespaceName, ",") {
		cache := c.client.GetConfigCache(v)
		cache.Range(func(key, value interface{}) bool {
			err = j.Set(dealPattern(gconv.String(key)), value)
			if err != nil {
				return false
			}
			return true
		})
	}
	if err == nil {
		c.value.Set(j)
	}
	return
}

// 兼容非properties格式中的key会被转小写的问题，干脆所有key都转成小写
func dealPattern(pattern string) string {
	return gstr.ToLower(pattern)
}
