package main

import (
	"fmt"
	"net/netip"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

func main() {
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Println(cCtx.App.Version)
	}

	app := &cli.App{
		Name:    "cidr-contains",
		Version: Version(),
		Usage:   "check whether a CIDR contains an IP address",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "cidr",
				Aliases: []string{"c"},
				Usage:   "CIDR (ex. 192.0.2.0/24, 2001:db8::/32)",
			},
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Usage:   "IP address (ex. 192.0.2.1, 2001:db8::1)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			return containsCommand(cCtx.String("cidr"), cCtx.String("address"))
		},
	}
	app.UsageText = fmt.Sprintf("%s [GLOBAL OPTIONS]", app.Name)

	if err := app.Run(os.Args); err != nil {
		cli.HandleExitCoder(err)
	}
}

const (
	exitCodeOK          = 0
	exitCodeNotContains = 1
	exitCodeUsageError  = 2
)

func containsCommand(cidrStr, addrStr string) error {
	cidr, err := netip.ParsePrefix(cidrStr)
	if err != nil {
		return cli.Exit(err.Error(), exitCodeUsageError)
	}
	addr, err := netip.ParseAddr(addrStr)
	if err != nil {
		return cli.Exit(err.Error(), exitCodeUsageError)
	}

	if !cidr.Contains(addr) {
		return cli.Exit("", exitCodeNotContains)
	}
	return nil
}

func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(devel)"
	}
	return info.Main.Version
}
