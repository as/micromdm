package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/as/micromdm/go4/version"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	if strings.ToLower(os.Args[1]) != "config" {
		checkForOldConfig()
	}
	var run func([]string) error
	switch strings.ToLower(os.Args[1]) {
	case "version", "-version":
		version.Print()
		return
	case "config":
		cmd := new(configCommand)
		run = cmd.Run
	case "get":
		cmd := new(getCommand)
		run = cmd.Run
	case "apply":
		cmd := new(applyCommand)
		run = cmd.Run
	case "remove":
		cmd := new(removeCommand)
		run = cmd.Run
	case "mdmcert":
		cmd := new(mdmcertCommand)
		run = cmd.Run
	default:
		usage()
		os.Exit(1)
	}

	if err := run(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func usage() error {
	helpText := `USAGE: mdmctl <COMMAND>

Available Commands:
	get
	apply
	config
	remove
	mdmcert
	version

Use micromdm <command> -h for additional usage of each command.
Example: mdmctl get devices
`
	fmt.Println(helpText)
	return nil
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
