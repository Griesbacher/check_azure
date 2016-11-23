package main

import (
	"github.com/griesbacher/check_azure/azureHttp"
	"github.com/griesbacher/check_azure/mode/subscription"
	"github.com/griesbacher/check_azure/mode/cvm"
	"github.com/griesbacher/check_x"
	"github.com/urfave/cli"
	"os"
	"time"
)

var (
	clientId string
	clientSecret string
	subscriptionId string
	tenantId string
	timeout int
	resourceGroup string
	name string
	grain string
	warning string
	critical string
	warn           *check_x.Threshold
	crit           *check_x.Threshold
)

func startTimeout() {
	if timeout != 0 {
		check_x.StartTimeout(time.Duration(timeout) * time.Second)
	}
}

func createConnector() azureHttp.AzureConnector {
	ac, err := azureHttp.NewAzureConnector(clientId, clientSecret, subscriptionId, tenantId)
	check_x.ExitOnError(err)
	return *ac
}

func parseThresholds() {
	var err error
	warn, err = check_x.NewThreshold(warning)
	check_x.ExitOnError(err)
	crit, err = check_x.NewThreshold(critical)
	check_x.ExitOnError(err)
}

func getReady(parse bool) azureHttp.AzureConnector {
	if parse {
		parseThresholds()
	}
	startTimeout()
	return createConnector()
}

func main() {
	app := cli.NewApp()
	app.Name = "check_azure"
	app.Usage = "Checks different azure stats\n   Copyright (c) 2016 Philip Griesbacher"
	app.Version = "0.0.1"
	VMFlagResourceGroup := cli.StringFlag{Name: "resourceGroup", Usage: "ResourceGroup of the VirtualMachine", Destination: &resourceGroup}
	VMFlagName := cli.StringFlag{Name: "name", Usage: "Name of the VirtualMachine", Destination: &name}
	VMFlagGrain := cli.StringFlag{Name: "grain", Usage: "Timegrain. Available: 5m,1h", Value: "5m", Destination: &grain}
	VMFlagWarning := cli.StringFlag{Name: "w", Usage: "warning: in the given timegrain", Destination: &warning, Value: "80"}
	VMFlagCritical := cli.StringFlag{Name: "c", Usage: "critical: in the given timegrain", Destination: &critical, Value: "90"}
	app.Commands = []cli.Command{
		{
			Name:    "mode",
			Aliases: []string{"m"},
			Usage:   "check mode",
			Subcommands: []cli.Command{
				{
					Name:    "subscriptions",
					Usage:   "Azure Subscriptions",
					Aliases: []string{"s"},
					Subcommands: []cli.Command{
						{
							Name: "show",
							Action: func(c *cli.Context) error {
								return subscription.Display(getReady(false))
							},
						},
					},
				}, {
					Name:    "classicVirtualMachines",
					Usage:   "classic VirtualMachines",
					Aliases: []string{"cvm"},
					Subcommands: []cli.Command{
						{
							Name:  "show",
							Usage: "List Mashines",
							Action: func(c *cli.Context) error {
								return cvm.ShowMachines(getReady(false), resourceGroup)
							},
							Flags: []cli.Flag{VMFlagResourceGroup},
						},
						{
							Name:  "cpu",
							Usage: "Percentage CPU",
							Action: func(c *cli.Context) error {
								return cvm.Cpu(getReady(true), resourceGroup, name, grain, warn, crit)
							},
							Flags: []cli.Flag{VMFlagResourceGroup, VMFlagName, VMFlagGrain, VMFlagWarning, VMFlagCritical},
						},
						{
							Name:  "network",
							Usage: "Network In/Out",
							Action: func(c *cli.Context) error {
								return cvm.Network(getReady(false), resourceGroup, name, grain, warning, critical)
							},
							Flags: []cli.Flag{VMFlagResourceGroup, VMFlagName, VMFlagGrain,
								cli.StringFlag{Name: "w", Usage: "warning: in,out in Bytes in the given timegrain", Destination: &warning, Value:"100000,100000"},
								cli.StringFlag{Name: "c", Usage: "critical: in,out in Bytes in the given timegrain", Destination: &critical, Value:"200000,200000"},
							},
						},
						{
							Name:  "disk",
							Usage: "Disk Read/Write",
							Action: func(c *cli.Context) error {
								return cvm.Disk(getReady(false), resourceGroup, name, grain, warning, critical)
							},
							Flags: []cli.Flag{VMFlagResourceGroup, VMFlagName, VMFlagGrain,
								cli.StringFlag{Name: "w", Usage: "warning: read,write in BytesPerSecond in the given timegrain", Destination: &warning, Value:"100,100"},
								cli.StringFlag{Name: "c", Usage: "critical: read,write in BytesPerSecond in the given timegrain", Destination: &critical, Value:"200,200"},
							},
						},
					},
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "clientId",
			Usage:       "Microsoft Client ID",
			Destination: &clientId,
		},
		cli.StringFlag{
			Name:        "clientSecret",
			Usage:       "Microsoft Client Secret",
			Destination: &clientSecret,
		},
		cli.StringFlag{
			Name:        "subscriptionId",
			Usage:       "Azure Subscription ID",
			Destination: &subscriptionId,
		},
		cli.StringFlag{
			Name:        "tenantId",
			Usage:       "Azure Tenant ID",
			Destination: &tenantId,
		},
		cli.IntFlag{
			Name:        "t",
			Usage:       "Seconds till check returns unknown, 0 to disable",
			Value:       10,
			Destination: &timeout,
		},
	}

	if err := app.Run(os.Args); err != nil {
		check_x.ExitOnError(err)
	}
}
