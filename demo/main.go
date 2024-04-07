package main

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	gcfg_apollo "github.com/yiqiang3344/gcfg-apollo"
	"time"
)

func main() {
	var (
		ctx = gctx.GetInitCtx()
	)
	adapter, err := gcfg_apollo.CreateAdapterApollo(ctx)
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}

	gcfg.Instance().SetAdapter(adapter)

	fmt.Printf("配置是否可用:%v\n", gcfg.Instance().Available(ctx))
	cfg, _ := gcfg.Instance().Data(ctx)
	fmt.Printf("所有配置数据:\n")
	g.Dump(cfg)
	for true {
		fmt.Printf("%v\n", gtime.Now().String())
		for _, _k := range []string{"test", "Test1.Test2"} {
			_v, _ := gcfg.Instance().Get(ctx, _k)
			fmt.Printf("%v:%v\n", _k, _v)
		}
		time.Sleep(5 * time.Second)
	}
}
