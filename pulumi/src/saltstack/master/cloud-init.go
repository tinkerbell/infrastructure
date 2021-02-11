package master

import (
	"fmt"

	"github.com/juju/juju/cloudconfig/cloudinit"
	"github.com/juju/packaging"
)

type BootstrapConfig struct {
	domain       string
	clientId     string
	clientSecret string
}

func cloudInitConfig(config *BootstrapConfig) string {
	c, err := cloudinit.New("focal")

	if err != nil {
		panic(err)
	}

	c.SetSystemUpdate(true)

	c.AddPackageSource(packaging.PackageSource{
		Name: "saltstack",
		URL:  "deb http://repo.saltstack.com/py3/ubuntu/20.04/amd64/3002 focal main",
		Key: `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v2

mQENBFOpvpgBCADkP656H41i8fpplEEB8IeLhugyC2rTEwwSclb8tQNYtUiGdna9
m38kb0OS2DDrEdtdQb2hWCnswxaAkUunb2qq18vd3dBvlnI+C4/xu5ksZZkRj+fW
tArNR18V+2jkwcG26m8AxIrT+m4M6/bgnSfHTBtT5adNfVcTHqiT1JtCbQcXmwVw
WbqS6v/LhcsBE//SHne4uBCK/GHxZHhQ5jz5h+3vWeV4gvxS3Xu6v1IlIpLDwUts
kT1DumfynYnnZmWTGc6SYyIFXTPJLtnoWDb9OBdWgZxXfHEcBsKGha+bXO+m2tHA
gNneN9i5f8oNxo5njrL8jkCckOpNpng18BKXABEBAAG0MlNhbHRTdGFjayBQYWNr
YWdpbmcgVGVhbSA8cGFja2FnaW5nQHNhbHRzdGFjay5jb20+iQE4BBMBAgAiBQJT
qb6YAhsDBgsJCAcDAgYVCAIJCgsEFgIDAQIeAQIXgAAKCRAOCKFJ3le/vhkqB/0Q
WzELZf4d87WApzolLG+zpsJKtt/ueXL1W1KA7JILhXB1uyvVORt8uA9FjmE083o1
yE66wCya7V8hjNn2lkLXboOUd1UTErlRg1GYbIt++VPscTxHxwpjDGxDB1/fiX2o
nK5SEpuj4IeIPJVE/uLNAwZyfX8DArLVJ5h8lknwiHlQLGlnOu9ulEAejwAKt9CU
4oYTszYM4xrbtjB/fR+mPnYh2fBoQO4d/NQiejIEyd9IEEMd/03AJQBuMux62tjA
/NwvQ9eqNgLw9NisFNHRWtP4jhAOsshv1WW+zPzu3ozoO+lLHixUIz7fqRk38q8Q
9oNR31KvrkSNrFbA3D89uQENBFOpvpgBCADJ79iH10AfAfpTBEQwa6vzUI3Eltqb
9aZ0xbZV8V/8pnuU7rqM7Z+nJgldibFk4gFG2bHCG1C5aEH/FmcOMvTKDhJSFQUx
uhgxttMArXm2c22OSy1hpsnVG68G32Nag/QFEJ++3hNnbyGZpHnPiYgej3FrerQJ
zv456wIsxRDMvJ1NZQB3twoCqwapC6FJE2hukSdWB5yCYpWlZJXBKzlYz/gwD/Fr
GL578WrLhKw3UvnJmlpqQaDKwmV2s7MsoZogC6wkHE92kGPG2GmoRD3ALjmCvN1E
PsIsQGnwpcXsRpYVCoW7e2nW4wUf7IkFZ94yOCmUq6WreWI4NggRcFC5ABEBAAGJ
AR8EGAECAAkFAlOpvpgCGwwACgkQDgihSd5Xv74/NggA08kEdBkiWWwJZUZEy7cK
WWcgjnRuOHd4rPeT+vQbOWGu6x4bxuVf9aTiYkf7ZjVF2lPn97EXOEGFWPZeZbH4
vdRFH9jMtP+rrLt6+3c9j0M8SIJYwBL1+CNpEC/BuHj/Ra/cmnG5ZNhYebm76h5f
T9iPW9fFww36FzFka4VPlvA4oB7ebBtquFg3sdQNU/MmTVV4jPFWXxh4oRDDR+8N
1bcPnbB11b5ary99F/mqr7RgQ+YFF0uKRE3SKa7a+6cIuHEZ7Za+zhPaQlzAOZlx
fuBmScum8uQTrEF5+Um5zkwC7EXTdH1co/+/V/fpOtxIg4XO4kcugZefVm5ERfVS
MA==
=dtMN
-----END PGP PUBLIC KEY BLOCK-----`})

	c.AddPackage("python3-pygit2")
	c.AddPackage("salt-master")
	c.AddPackage("salt-minion")

	c.AddRunTextFile("/etc/salt/master.d/master.conf", `autosign_grains_dir: /etc/salt/autosign-grains
fileserver_backend:
  - roots
  - gitfs

gitfs_remotes:
  - https://github.com/tinkerbell/infrastructure:
      - root: saltstack/states
      - base: main
      - update_interval: 120
  - https://github.com/saltstack-formulas/fail2ban-formula:
      - base: master

pillar_roots:
  base:
    - /srv/pillar

ext_pillar:
  - git:
      - main https://github.com/tinkerbell/infrastructure:
          - root: saltstack/pillars
          - env: main`, 0644)

	c.AddRunTextFile("/etc/salt/minion.d/minion.conf", `autosign_grains:
- role

grains:
  role: master`, 0644)

	c.AddRunCmd("mkdir -p /srv/pillar")
	c.AddRunTextFile("/srv/pillar/top.sls", `base:
  'G@role:master':
    - teleport
`, 0644)
	c.AddRunTextFile("/srv/pillar/teleport.sls", fmt.Sprintf("teleport:\n  domain: %s\n  clientId: %s\n  clientSecret: %s\n", config.domain, config.clientId, config.clientSecret), 0644)

	c.AddRunCmd("PRIVATE_IP=$(curl -s https://metadata.platformequinix.com/metadata | jq -r '.network.addresses | map(select(.public==false)) | first | .address')")
	c.AddRunCmd("echo interface: ${PRIVATE_IP} > /etc/salt/master.d/private-interface.conf")
	c.AddRunCmd("echo master: ${PRIVATE_IP} > /etc/salt/minion.d/master.conf")
	c.AddRunCmd("mkdir -p /etc/salt/autosign-grains/")
	c.AddRunCmd("echo master > /etc/salt/autosign-grains/role")
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
