package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type depDevicesTableOutput struct{ w *tabwriter.Writer }

func (out *depDevicesTableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "SerialNumber\tModel\tProfileStatus\tProfileUUID\n")
}

func (out *depDevicesTableOutput) BasicFooter() {
	out.w.Flush()
}

func (cmd *getCommand) getDEPDevices(args []string) error {
	flagset := flag.NewFlagSet("dep-devices", flag.ExitOnError)
	flSerials := flagset.String("serials", "", "comma separated list of device serials")
	flagset.Usage = usageFor(flagset, "mdmctl get dep-devices [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	if *flSerials == "" {
		flagset.Usage()
		return errors.New("bad input: must provide a comma separated list of DEP serials")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	out := &depDevicesTableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()
	ctx := context.Background()
	serials := strings.Split(*flSerials, ",")
	resp, err := cmd.depsvc.GetDeviceDetails(ctx, serials)
	if err != nil {
		return err
	}
	for _, d := range resp.Devices {
		fmt.Fprintf(out.w, "%s\t%s\t%s\t%s\n", d.SerialNumber, d.Model, d.ProfileStatus, d.ProfileUUID)
	}
	return nil
}
