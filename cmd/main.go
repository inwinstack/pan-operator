package main

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/inwinstack/pa-operator/pkg/operator"
	"github.com/inwinstack/pa-operator/pkg/pautil"
	"github.com/inwinstack/pa-operator/pkg/version"
	flag "github.com/spf13/pflag"
)

var (
	kubeconfig string
	host       string
	username   string
	password   string
	namespaces []string
	retry      int
	commitTime int
	ver        bool
)

func parserFlags() {
	flag.StringVarP(&kubeconfig, "kubeconfig", "", "", "Absolute path to the kubeconfig file.")
	flag.StringVarP(&host, "pa-host", "", "", "Palo Alto API host address.")
	flag.StringVarP(&username, "pa-username", "", "", "Palo Alto API username.")
	flag.StringVarP(&password, "pa-password", "", "", "Palo Alto API password.")
	flag.StringSliceVarP(&namespaces, "ignore-namespaces", "", nil, "Set ignore namespaces for Kubernetes service.")
	flag.IntVarP(&retry, "retry", "", 5, "Number of retry for PA failed job.")
	flag.IntVarP(&commitTime, "commit-wait-time", "", 2, "The length of time to wait next PA commit.")
	flag.BoolVarP(&ver, "version", "", false, "Display the version.")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}

func main() {
	defer glog.Flush()
	parserFlags()

	glog.Infof("Starting PA operator...")

	if ver {
		fmt.Fprintf(os.Stdout, "%s\n", version.GetVersion())
		os.Exit(0)
	}

	f := &operator.Flag{
		Kubeconfig:       kubeconfig,
		IgnoreNamespaces: namespaces,
		Retry:            retry,
		CommitWaitTime:   commitTime,
		PaloAlto: &pautil.Flag{
			Host:     host,
			Username: username,
			Password: password,
		},
	}

	op := operator.NewMainOperator(f)
	if err := op.Initialize(); err != nil {
		glog.Fatalf("Error initing operator instance: %v.", err)
	}

	if err := op.Run(); err != nil {
		glog.Fatalf("Error serving operator instance: %s.", err)
	}
}
