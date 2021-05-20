package master

import (
	"fmt"

	"github.com/juju/juju/cloudconfig/cloudinit"
)

type BootstrapConfig struct {
	teleportDomain       string
	teleportClientID     string
	teleportClientSecret string
	teleportPeerToken    string
	githubUsername       string
	githubAccessToken    string
	awsAccessKeyID       string
	awsSecretAccessKey   string
	awsBucketName        string
	awsBucketLocation    string
}

func cloudInitConfig(config *BootstrapConfig) string {
	c, err := cloudinit.New("focal")

	if err != nil {
		panic(err)
	}

	c.SetSystemUpdate(true)

	c.AddRunCmd("curl -fsSL https://bootstrap.saltproject.io -o install_salt.sh")
	c.AddRunCmd("sh install_salt.sh -P -M -x python3")

	c.AddRunTextFile("/etc/salt/master.d/master.conf", `autosign_grains_dir: /etc/salt/autosign-grains
fileserver_backend:
  - roots
  - gitfs

gitfs_remotes:
  - https://github.com/tinkerbell/infrastructure:
      - root: saltstack/salt
      - base: main
      - update_interval: 120
  - https://github.com/saltstack-formulas/fail2ban-formula:
      - base: master

pillar_roots:
  base:
    - /srv/pillar

ext_pillar:
  - http_json:
      url: https://metadata.platformequinix.com/metadata
  - git:
      - main https://github.com/tinkerbell/infrastructure:
          - root: saltstack/pillar
          - env: main`, 0644)

	c.AddRunTextFile("/etc/salt/minion.d/minion.conf", `autosign_grains:
- role

startup_states: highstate
schedule:
  highstate:
    function: state.highstate
    minutes: 15

grains:
  role: master`, 0644)

	c.AddRunCmd("PRIVATE_IP=$(curl -s https://metadata.platformequinix.com/metadata | jq -r '.network.addresses | map(select(.public==false)) | first | .address')")

	c.AddRunCmd("mkdir -p /srv/pillar")
	c.AddRunCmd("echo salt_master_private_ipv4: ${$PRIVATE_IP} > /srv/pillar/base.sls")
	c.AddRunTextFile("/srv/pillar/top.sls", `base:
  '*':
    - base
	- teleport.node

  'G@role:master':
    - aws
    - teleport
	
  'G@role:github-action-runner':
    - github
`, 0644)

	c.AddRunCmd("mkdir -p /srv/pillar/teleport")
	c.AddRunTextFile("/srv/pillar/teleport/init.sls", fmt.Sprintf("teleport:\n  domain: %s\n  clientId: %s\n  clientSecret: %s\n", config.teleportDomain, config.teleportClientID, config.teleportClientSecret), 0400)
	c.AddRunTextFile("/srv/pillar/teleport/node.sls", fmt.Sprintf("teleport:\n  peerToken: %s\n", config.teleportPeerToken), 0400)
	c.AddRunTextFile("/srv/pillar/github.sls", fmt.Sprintf("github:\n  username: %s\n  accessToken: %s\n", config.githubUsername, config.githubAccessToken), 0400)
	c.AddRunTextFile("/srv/pillar/aws.sls", fmt.Sprintf("s3.keyid: %s\ns3.key: %s\ns3.location: %s\ns3.bucketName: %s\n", config.awsAccessKeyID, config.awsSecretAccessKey, config.awsBucketLocation, config.awsBucketName), 0400)

	c.AddRunCmd("echo interface: ${PRIVATE_IP} > /etc/salt/master.d/private-interface.conf")
	c.AddRunCmd("echo master: ${PRIVATE_IP} > /etc/salt/minion.d/master.conf")
	c.AddRunCmd("mkdir -p /etc/salt/autosign-grains/")
	c.AddRunCmd("echo -e \"master\ngithub-action-runner\n\" > /etc/salt/autosign-grains/role")
	c.AddRunCmd("systemctl daemon-reload")
	c.AddRunCmd("systemctl enable salt-master.service")
	c.AddRunCmd("systemctl restart --no-block salt-master.service")
	c.AddRunCmd("systemctl enable salt-minion.service")
	c.AddRunCmd("systemctl restart --no-block salt-minion.service")

	script, err := c.RenderScript()

	if err != nil {
		panic(err)
	}

	return script

}
