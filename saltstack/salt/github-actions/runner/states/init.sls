{% if grains['gha_runner_states'] is defined %}
{% for state in grains['gha_runner_states'] %}
include:
- github-actions.runner.states.{{ state }}
{% endfor %}
{% endif %}