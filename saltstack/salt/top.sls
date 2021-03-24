base:
  '*':
    - equinix-metal
    - fail2ban

  'G@role:master':
    - teleport

  'G@role:github-action-runner':
    - github-actions.runner
    - github-actions.runner.register
    - github-actions.runner.states
