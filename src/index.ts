import { Facility, Plan } from "@pulumi/equinix-metal";
import * as pulumi from "@pulumi/pulumi";

import { createSaltMaster } from "./saltstack";

const saltMaster = createSaltMaster({
  name: "salt-master",
  facilities: [Facility.AM6],
  plan: Plan.C3MediumX86,
});
