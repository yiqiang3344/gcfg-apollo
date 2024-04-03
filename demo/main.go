package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	gcfg_apollo "github.com/yiqiang3344/gcfg-apollo"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "gcfg apollo test",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			//使用apollo的配置
			adapter, err := gcfg_apollo.CreateAdapterApollo(ctx)
			if err != nil {
				return err
			}
			gcfg.Instance().SetAdapter(adapter)

			fmt.Printf("配置是否可用:%v\n", gcfg.Instance().Available(ctx))
			v, _ := gcfg.Instance().Get(ctx, "projectName")
			fmt.Printf("application中的projectName对应值:%v\n", v)
			cfg, _ := gcfg.Instance().Data(ctx)
			fmt.Printf("所有配置数据:\n")
			g.Dump(cfg)

			return nil
		},
	}
)

func main() {
	Main.Run(gctx.New())
}
