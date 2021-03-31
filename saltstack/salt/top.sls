base:
  '*':
    - equinix-metal
    - fail2ban

  'G@role:master':
    - teleport

  'not G@role:master:
    - teleport.node

  'G@role:github-action-runner':
    - go
    - github-actions.runner
    - github-actions.runner.register
    - github-actions.runner.states
