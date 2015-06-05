package servercommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var reboot = cli.Command{
	Name:        "reboot",
	Usage:       fmt.Sprintf("%s %s reboot [--id <serverID> | --name <serverName>] [--soft | --hard] [optional flags]", util.Name, commandPrefix),
	Description: "Reboots an existing server",
	Action:      commandReboot,
	Flags:       util.CommandFlags(flagsReboot),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsReboot))
	},
}

func flagsReboot() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "soft",
			Usage: "[optional; required if 'hard' is not provided] Ask the OS to restart under its own procedures.",
		},
		cli.BoolFlag{
			Name:  "hard",
			Usage: "[optional; required if 'soft' is not provided] Physically cut power to the machine and then restore it after a brief while.",
		},
	}
}

func commandReboot(c *cli.Context) {
	util.CheckArgNum(c, 0)

	var how osServers.RebootMethod
	if c.IsSet("soft") {
		how = osServers.OSReboot
	}
	if c.IsSet("hard") {
		how = osServers.PowerCycle
	}

	if how == "" {
		fmt.Printf("Missing flag: One of either --soft or --hard must be provided.")
		os.Exit(1)
	}

	client := auth.NewClient("compute")
	serverID := idOrName(c, client)
	err := servers.Reboot(client, serverID, how).ExtractErr()
	if err != nil {
		fmt.Printf("Error rebooting server (%s): %s\n", serverID, err)
		os.Exit(1)
	}
}