# Tinkerbell Infrastructure

## Pulumi

Coming soon

## SaltStack

### Prereqs

You need to run `nix-shell` to get the required dependencies.

### Tests

From within the `nix-shell`, you can `cd saltstack` and run `kitchen test`.

If you're testing a new state, please add it to `kitchen.yml` to ensure it runs.
