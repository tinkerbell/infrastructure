package githubactions

import (
	"fmt"

	"github.com/juju/juju/cloudconfig/cloudinit"
)

type MinionConfig struct {
	MasterIp string
}

func cloudInitConfig(config *MinionConfig) string {
	c, err := cloudinit.New("focal")

	if err != nil {
		panic(err)
	}

	c.SetSystemUpdate(true)

	c.AddRunCmd("curl -fsSL https://bootstrap.saltproject.io -o install_salt.sh")
	c.AddRunCmd("sh install_salt.sh -P -x python3")
	c.AddRunTextFile("/etc/salt/minion.d/minion.conf", `autosign_grains:
- role

startup_states: highstate
schedule:
  highstate:
    function: state.highstate
    minutes: 15

grains:
  role: github-action-runner`, 0644)

	c.AddRunCmd(fmt.Sprintf("echo master: %s > /etc/salt/minion.d/master.conf", config.MasterIp))
	c.AddRunCmd("systemctl daemon-reload")
	c.AddRunCmd("systemctl enable salt-minion.service")
	c.AddRunCmd("systemctl restart --no-block salt-minion.service")

	script, err := c.RenderScript()

	if err != nil {
		panic(err)
	}

	return script

}
