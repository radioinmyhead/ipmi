package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/radioinmyhead/ipmi"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "goipmi"
	app.Version = "0.1"
	app.Usage = "Control the IPMI setting"

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
		}, {
			Name:  "fan-status",
			Usage: "get fan status",
			Action: func(c *cli.Context) error {
				bmc, err := ipmi.GetLocalIPMI()
				if err != nil {
					return err
				}
				return bmc.GetFanSpeed()
			},
		}, {
			Name:  "fan-rpm",
			Usage: "get fan rpm",
			Action: func(c *cli.Context) error {
				bmc, err := ipmi.GetLocalIPMI()
				if err != nil {
					return err
				}
				return bmc.GetFanRPM()
			},
		}, {
			Name:  "set-conf",
			Usage: "set conf with conf-file",
			Action: func(c *cli.Context) error {
				type Conf struct {
					Fanset int `json:"fan-set"`
				}
				cfile := c.Args().First()
				if cfile == "" {
					cfile = "/etc/megvii/ipmi.json"
				}
				raw, err := ioutil.ReadFile(cfile)
				if err != nil {
					return err
				}
				var conf Conf
				if err := json.Unmarshal(raw, &conf); err != nil {
					return err
				}
				logrus.Info(conf)
				if conf.Fanset != 0 {
					bmc, err := ipmi.GetLocalIPMI()
					if err != nil {
						return err
					}
					if err = bmc.CheckSpeed(conf.Fanset); err != nil {
						return err
					}
					if err = bmc.SetFanSpeed(conf.Fanset); err != nil {
						return err
					}
				}
				return nil
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
