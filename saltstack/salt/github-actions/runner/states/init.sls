{% if grains['gha_runner_states'] is defined %}
include:
{% for state in grains['gha_runner_states'] %}
- github-actions.runner.states.{{ state }}
{% endfor %}
{% endif %}