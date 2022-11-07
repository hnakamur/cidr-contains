package main

import (
	"fmt"
	"net/netip"
	"os"
	"runtime/debug"
	"strings"

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
				Name:     "cidr",
				Aliases:  []string{"c"},
				Value:    &NetPrefixValue{},
				Required: true,
				Usage:    "CIDR (ex. 192.0.2.0/24, 2001:db8::/32)",
			},
			&cli.GenericFlag{
				Name:     "address",
				Aliases:  []string{"a"},
				Value:    &NetAddrParser{},
				Required: true,
				Usage:    "IP address (ex. 192.0.2.1, 2001:db8::1)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			cidr := cCtx.Generic("cidr").(*NetPrefixValue).prefix
			addr := cCtx.Generic("address").(*NetAddrParser).addr
			return containsCommand(cidr, addr)
		},
	}
	app.UsageText = fmt.Sprintf("%s [GLOBAL OPTIONS]", app.Name)
	app.OnUsageError = func(cCtx *cli.Context, err error, isSubcommand bool) error {
		return cli.Exit(err.Error(), exitCodeUsageError)
	}

	if err := app.Run(os.Args); err != nil {
		cli.HandleExitCoder(err)
		fmt.Fprintf(app.ErrWriter, "\nError: %s\n", err)
		if strings.HasPrefix(err.Error(), "Required flag") {
			os.Exit(exitCodeUsageError)
		}
	}
}

func containsCommand(cidr netip.Prefix, addr netip.Addr) error {
	if !cidr.Contains(addr) {
		return cli.Exit("", exitCodeNotContains)
	}
	return nil
}

type NetPrefixValue struct {
	prefix     netip.Prefix
	hasBeenSet bool
}

func (p *NetPrefixValue) Set(value string) error {
	prefix, err := netip.ParsePrefix(value)
	if err != nil {
		return err
	}
	*p = NetPrefixValue{
		prefix:     prefix,
		hasBeenSet: true,
	}
	return nil
}

func (p *NetPrefixValue) String() string {
	if p.hasBeenSet {
		return p.prefix.String()
	}
	return ""
}

func (p *NetPrefixValue) Get() any {
	if p.hasBeenSet {
		return p.prefix
	}
	return nil
}

type NetAddrParser struct {
	addr       netip.Addr
	hasBeenSet bool
}

func (p *NetAddrParser) Set(value string) error {
	addr, err := netip.ParseAddr(value)
	if err != nil {
		return err
	}
	*p = NetAddrParser{
		addr:       addr,
		hasBeenSet: true,
	}
	return nil
}

func (p *NetAddrParser) String() string {
	if p.hasBeenSet {
		return p.addr.String()
	}
	return ""
}

func (p *NetAddrParser) Get() any {
	if p.hasBeenSet {
		return p.addr
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
