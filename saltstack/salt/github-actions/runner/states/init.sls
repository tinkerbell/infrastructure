{% if grains['gha_runner_states'] is defined %}
include:
- github-actions.runner.states.nix
{% for state in grains['gha_runner_states'] %}
- github-actions.runner.states.{{ state }}
{% endfor %}
{% endif %}
