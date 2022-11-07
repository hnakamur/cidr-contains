package main

import (
	"fmt"
	"net/netip"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

const (
	exitCodeOK          = 0
	exitCodeNotContains = 1
	exitCodeUsageError  = 2
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
			&cli.GenericFlag{
				Name:    "cidr",
				Aliases: []string{"c"},
				Value:   &NetPrefix{},
				Usage:   "CIDR (ex. 192.0.2.0/24, 2001:db8::/32)",
			},
			&cli.GenericFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Value:   &NetAddr{},
				Usage:   "IP address (ex. 192.0.2.1, 2001:db8::1)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			cidr := netip.Prefix(*cCtx.Generic("cidr").(*NetPrefix))
			addr := netip.Addr(*cCtx.Generic("address").(*NetAddr))
			return containsCommand(cidr, addr)
		},
	}
	app.UsageText = fmt.Sprintf("%s [GLOBAL OPTIONS]", app.Name)
	app.OnUsageError = func(cCtx *cli.Context, err error, isSubcommand bool) error {
		return cli.Exit(err.Error(), exitCodeUsageError)
	}

	if err := app.Run(os.Args); err != nil {
		cli.HandleExitCoder(err)
	}
}

func containsCommand(cidr netip.Prefix, addr netip.Addr) error {
	if !cidr.Contains(addr) {
		return cli.Exit("", exitCodeNotContains)
	}
	return nil
}

type NetPrefix netip.Prefix

func (p *NetPrefix) Set(value string) error {
	parsed, err := netip.ParsePrefix(value)
	if err != nil {
		return err
	}
	*p = NetPrefix(parsed)
	return nil
}

func (p *NetPrefix) String() string {
	return netip.Prefix(*p).String()
}

type NetAddr netip.Addr

func (p *NetAddr) Set(value string) error {
	parsed, err := netip.ParseAddr(value)
	if err != nil {
		return err
	}
	*p = NetAddr(parsed)
	return nil
}

func (p *NetAddr) String() string {
	return netip.Addr(*p).String()
}

func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(devel)"
	}
	return info.Main.Version
}
