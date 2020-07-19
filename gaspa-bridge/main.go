// Command gaspa-bridge runs a simple bridge service
package main

import "github.com/urfave/cli/v2"

func main() {
	app := cli.NewApp()
	app.Name = "Gaspa-Bridge"
	app.Usage = "Gaspa Standalone Bridge Service"
	app.UsageText = "gaspa-bridge [CONFIGURE_FILE]"
	app.Description = `Gaspa bridge is used when several nodes cannot communicate
	with any methods as the network environment is too complex.
	It runs a simple transparent TCP proxy to let traffics pass
	by. It needs to be deployed with a public network IP that's
	visible to all of the nodes.`
	app.HideHelpCommand = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "configure",
			Aliases:     []string{"c"},
			DefaultText: "/etc/gaspa/bridge.yml",
			Usage:       "Configure file path",
		},
	}
	app.Action = func(ctx *cli.Context) error {

		return nil
	}
	app.RunAndExitOnError()
}
