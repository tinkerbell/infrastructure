import * as pulumi from "@pulumi/pulumi";
import {
  BillingCycle,
  Device,
  Facility,
  OperatingSystem,
  Plan,
} from "@pulumi/equinix-metal";
import { readFileSync } from "fs";
import * as yaml from "js-yaml";
import { join } from "path";

import { KEY } from "../apt-key";

interface SaltMasterConfig {
  name: string;
  facilities: Facility[];
  plan: Plan;
}

interface SaltMaster {
  device: Device;
}

interface WriteFile {
  content: string;
  path: string;
  permissions?: string;
}

const writeFile = (filename: string, path: string): WriteFile => ({
  content: readFileSync(join(__dirname, "files", filename)).toString(),
  path,
});

export const createSaltMaster = (config: SaltMasterConfig): SaltMaster => {
  const projectId = new pulumi.Config("equinix-metal").require("projectId");

  const cloudConfig = yaml.dump({
    apt: {
      sources: {
        salt: {
          source:
            "deb http://repo.saltstack.com/py3/ubuntu/20.04/amd64/3002 focal main",
          key: KEY,
          filename: "saltstack.list",
        },
      },
    },
    packages: ["python3-pygit2", "salt-master", "salt-minion"],
    write_files: [
      writeFile("master.d.yaml", "/etc/salt/master.d/master.conf"),
      writeFile("minion.d.yaml", "/etc/salt/minion.d/minion.conf"),
    ],
    runcmd: [
      "PRIVATE_IP=$(curl -s https://metadata.platformequinix.com/metadata | jq -r '.network.addresses | map(select(.public==false)) | first | .address')",
      "echo interface: ${PRIVATE_IP} > /etc/salt/master.d/private-interface.conf",
      "echo master: ${PRIVATE_IP} > /etc/salt/minion.d/master.conf",
      "mkdir -p /etc/salt/autosign-grains/",
      "echo master > /etc/salt/autosign-grains/role",
      "systemctl daemon-reload",
      "systemctl enable salt-master.service",
      "systemctl restart --no-block salt-master.service",
      "systemctl enable salt-minion.service",
      "systemctl restart --no-block salt-minion.service",
    ],
  });

  return {
    device: new Device(config.name, {
      billingCycle: BillingCycle.Hourly,
      facilities: config.facilities,
      hostname: config.name,
      operatingSystem: OperatingSystem.Ubuntu2010,
      plan: config.plan,
      projectId,
      tags: ["role:salt-master"],
      userData: `#cloud-config\n${cloudConfig}\n`,
    }),
  };
};
