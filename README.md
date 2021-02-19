# Tinkerbell Infrastructure

## Pulumi

Coming soon

## SaltStack

### Prereqs

You need to run `nix-shell` to get the required dependencies.

### Tests

From within the `nix-shell`, you can `cd saltstack` run `bundle install --binstubs` and then run `bundle exec bin/kitchen test`.

If you're testing a new state, please add it to `kitchen.yml` to ensure it runs.
