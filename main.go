package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hostwithquantum/ansible-openstack-inventory/auth"
	"github.com/hostwithquantum/ansible-openstack-inventory/file"
	"github.com/hostwithquantum/ansible-openstack-inventory/inventory"
	"github.com/hostwithquantum/ansible-openstack-inventory/server"
	"github.com/urfave/cli/v2"
	"gopkg.in/ini.v1"
)

var version string
var defaultGroup = "all"

// FIXME: move into config.ini
var defaultVars = map[string]string{
	"ansible_ssh_user":           "core",
	"ansible_ssh_common_args":    "-F customers/files/quantum/ssh_config",
	"ansible_python_interpreter": "/opt/python/bin/python",
}

func main() {
	app := &cli.App{
		Name:    "ansible-openstack-inventory",
		Usage:   "A cli tool for dynamic inventories for Planetary Quantum",
		Version: version,
		Authors: []*cli.Author{
			{
				Name:  "Till Klampaeckel",
				Email: "till@planetary-quantum.com",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "list",
				Usage: "List the repository",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "host",
				Usage: "List an individual host",
			},
			&cli.StringFlag{
				Name:    "customer",
				Usage:   "The customer to pull nodes for",
				EnvVars: []string{"QUANTUM_CUSTOMER"},
			},
			&cli.StringFlag{
				Name:    "config",
				Usage:   "Settings for groups, etc. for --list",
				Value:   "config.ini",
				EnvVars: []string{"QUANTUM_INVENTORY_CONFIG"},
			},
			&cli.StringFlag{
				Name:    "load-group-vars",
				Usage:   "Path to ./inventory/customer/group_vars",
				Value:   "",
				EnvVars: []string{"QUANTUM_INVENTORY_VARS_PATH"},
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("list") && c.String("host") != "" {
				log.Fatal("Can only use one of `--list` or `--host node`.")
			}

			if c.String("host") == "" && !c.Bool("list") {
				log.Fatal("No command provided.")
			}

			provider, err := auth.Authenticate()
			if err != nil {
				log.Fatal(err)
			}

			cfg, err := ini.Load(c.String("config"))
			if err != nil {
				log.Fatal(err)
			}
			var accessNetwork = cfg.Section("").Key("network").String()

			api := server.NewAPI(accessNetwork, provider)

			if c.String("host") != "" {
				s := api.GetByNode(c.String("host"))

				hostVars := make(map[string]string)
				hostVars["ansible_host"] = s.IPAddress
				hostVars["floating_ip"] = s.FloatingIP

				json, err := json.Marshal(hostVars)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(string(json))
				os.Exit(0)
			}

			var childrenGroups = cfg.Section("all").Key("children").Strings(",")

			customer := c.String("customer")
			if customer == "" {
				log.Fatal("No customer env variable")
			}

			allServers := api.GetByCustomer(customer)

			inventory := inventory.NewInventory(customer, append(childrenGroups, defaultGroup))

			for defaultVar, defaultValue := range defaultVars {
				inventory.AddVarToGroup(defaultGroup, defaultVar, defaultValue)
			}

			for _, server := range allServers {
				inventory.AddHostToGroup(server.Name, defaultGroup)
				inventory.AddHostVar("ansible_host", server.IPAddress, server.Name)

				if server.FloatingIP != "" {
					inventory.AddHostVar("floating_ip", server.FloatingIP, server.Name)
				}

				for _, g := range childrenGroups {
					inventory.AddHostToGroup(server.Name, g)
				}
			}

			groupVarsFile := file.NewGroupVarsFile(c.String("load-group-vars"))

			for _, group := range append(childrenGroups, defaultGroup) {
				sec, err := cfg.GetSection(group)
				if err != nil {
					continue
				}

				inventory.AddChildrenToGroup(sec.Key("children").Strings(","), group)

				if c.String("load-group-vars") != "" {
					groupFileYaml, err := groupVarsFile.HandleGroup(group)
					if err != nil {
						//log.Println(err)
						continue
					}

					for varKey, varValue := range groupFileYaml {
						inventory.AddVarToGroup(group, varKey, varValue)
					}
				}
			}

			fmt.Println(inventory.ReturnJSONInventory())
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
