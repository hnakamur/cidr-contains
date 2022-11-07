# cidr-contains

## Usage

```
$ cidr-contains -h
NAME:
   cidr-contains - check whether a CIDR contains an IP address

USAGE:
   cidr-contains [GLOBAL OPTIONS]

VERSION:
   v0.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cidr value, -c value     CIDR (ex. 192.0.2.0/24, 2001:db8::/32) (default: invalid Prefix)
   --address value, -a value  IP address (ex. 192.0.2.1, 2001:db8::1) (default: invalid IP)
   --help, -h                 show help (default: false)
   --version, -v              print the version (default: false)
```

## How to install

Download and install [Go](https://go.dev/) and run the following command:

```
go install github.com/hnakamur/cidr-contains@latest
```
