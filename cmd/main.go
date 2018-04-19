package main

import (
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/radioinmyhead/ipmi"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "goipmi"
	app.Usage = "Control the IPMI setting"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "fan-set",
			Usage: "set fan speed",
			Action: func(c *cli.Context) error {
				bmc, err := ipmi.GetLocalIPMI()
				if err != nil {
					return err
				}
				s, err := strconv.Atoi(c.Args().First())
				if err != nil {
					return err
				}
				if err = bmc.CheckSpeed(s); err != nil {
					return err
				}
				return bmc.SetFanSpeed(s)
			},
		},
		{
			Name:  "fan-status",
			Usage: "get fan status",
			Action: func(c *cli.Context) error {
				bmc, err := ipmi.GetLocalIPMI()
				if err != nil {
					return err
				}
				return bmc.GetFanSpeed()
			},
		},
		{
			Name:  "fan-rpm",
			Usage: "get fan rpm",
			Action: func(c *cli.Context) error {
				bmc, err := ipmi.GetLocalIPMI()
				if err != nil {
					return err
				}
				return bmc.GetFanRPM()
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
