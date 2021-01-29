import * as pulumi from "@pulumi/pulumi";
import {
  BillingCycle,
  Device,
  Facility,
  OperatingSystem,
  Plan,
} from "@pulumi/equinix-metal";
import * as yaml from "js-yaml";

import { KEY } from "./apt-key";

interface SaltMasterConfig {
  name: string;
  facilities: Facility[];
  plan: Plan;
}

interface SaltMaster {
  device: Device;
}

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
    packages: ["salt-master"],
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
