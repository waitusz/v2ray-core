package api

import (
	"fmt"

	handlerService "v2ray.com/core/app/proxyman/command"
	"v2ray.com/core/common/cmdarg"
	"v2ray.com/core/infra/conf/serial"
	"v2ray.com/core/main/commands/base"
)

var cmdRemoveOutbounds = &base.Command{
	CustomFlags: true,
	UsageLine:   "{{.Exec}} api rmo [--server=127.0.0.1:8080] <json_file|tag> [json_file] [tag]...",
	Short:       "Remove outbounds",
	Long: `
Remove outbounds from V2Ray.

Arguments:

	-s, -server 
		The API server address. Default 127.0.0.1:8080

	-t, -timeout
		Timeout seconds to call API. Default 3

Example:

    {{.Exec}} {{.LongName}} --server=127.0.0.1:8080 c1.json "tag name"
`,
	Run: executeRemoveOutbounds,
}

func executeRemoveOutbounds(cmd *base.Command, args []string) {
	setSharedFlags(cmd)
	cmd.Flag.Parse(args)
	unnamedArgs := cmd.Flag.Args()
	if len(unnamedArgs) == 0 {
		fmt.Println("reading from stdin:")
		unnamedArgs = []string{"stdin:"}
	}

	tags := make([]string, 0)
	for _, arg := range unnamedArgs {
		if r, err := cmdarg.LoadArg(arg); err == nil {
			conf, err := serial.DecodeJSONConfig(r)
			if err != nil {
				base.Fatalf("failed to decode %s: %s", arg, err)
			}
			outs := conf.OutboundConfigs
			for _, o := range outs {
				tags = append(tags, o.Tag)
			}
		} else {
			// take request as tag
			tags = append(tags, arg)
		}
	}

	if len(tags) == 0 {
		base.Fatalf("no outbound to remove")
	}

	conn, ctx, close := dialAPIServer()
	defer close()

	client := handlerService.NewHandlerServiceClient(conn)
	for _, tag := range tags {
		fmt.Println("removing:", tag)
		r := &handlerService.RemoveOutboundRequest{
			Tag: tag,
		}
		resp, err := client.RemoveOutbound(ctx, r)
		if err != nil {
			base.Fatalf("failed to remove outbound: %s", err)
		}
		showResponese(resp)
	}
}