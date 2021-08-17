package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hostwithquantum/ansible-openstack-inventory/auth"
	"github.com/hostwithquantum/ansible-openstack-inventory/file"
	"github.com/hostwithquantum/ansible-openstack-inventory/fip"
	"github.com/hostwithquantum/ansible-openstack-inventory/host"
	"github.com/hostwithquantum/ansible-openstack-inventory/inventory"
	"github.com/hostwithquantum/ansible-openstack-inventory/lbaas"
	"github.com/hostwithquantum/ansible-openstack-inventory/presenter"
	"github.com/hostwithquantum/ansible-openstack-inventory/response"
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
				return errors.New("Can only use one of `--list` or `--host node`.")
			}

			if c.String("host") == "" && !c.Bool("list") {
				return errors.New("No command provided.")
			}

			customer := c.String("customer")
			if customer == "" {
				return errors.New("No customer env variable")
			}

			provider, err := auth.Authenticate()
			if err != nil {
				return err
			}

			cfg, err := ini.Load(c.String("config"))
			if err != nil {
				return err
			}
			var accessNetwork = cfg.Section("").Key("network").String()

			fip := fip.NewFIP(accessNetwork, provider)
			lb := lbaas.NewAPI(customer, provider)
			compute := server.NewAPI(customer, accessNetwork, provider)

			p := presenter.Presenter{
				FIPs: fip.GetIps(),
			}

			if c.String("host") != "" {
				server, err := compute.GetByNode(c.String("host"))
				if err != nil {
					return err
				}

				server = p.AddFipToNode(server)
				if lb.HasLB() {
					haproxy, err := lb.GetByName()
					if err != nil {
						return err
					}
					server = p.AddLoadbalancerToNode(server, haproxy)
				}

				json, err := json.Marshal(host.Build(server))
				if err != nil {
					return err
				}

				fmt.Println(string(json))
				return nil
			}

			var childrenGroups = cfg.Section("all").Key("children").Strings(",")

			allServers, err := compute.GetByCustomer(customer)
			if err != nil {
				return err
			}

			if len(allServers) == 0 {
				// return early and avoid odd warnings when invoked via Ansible
				fmt.Println(response.BuildEmptyRepository(nil))
				return nil
			}

			allServers = p.AddFipsToNodes(allServers)
			if lb.HasLB() {
				haproxy, err := lb.GetByName()
				if err != nil {
					return err
				}

				var ready []server.AnsibleServer
				for _, server := range allServers {
					server = p.AddLoadbalancerToNode(server, haproxy)
					ready = append(ready, server)
				}

				allServers = ready
			}

			inventory := inventory.NewInventory(customer, append(childrenGroups, defaultGroup))

			for defaultVar, defaultValue := range defaultVars {
				inventory.AddVarToGroup(defaultGroup, defaultVar, defaultValue)
			}

			inventory.BuildServers(allServers, defaultGroup)
			inventory.BuildServerGroups(allServers, childrenGroups)

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
		fmt.Print(response.BuildEmptyRepository(err))
		os.Exit(1)
	}
	os.Exit(0)
}
