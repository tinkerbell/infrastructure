{% set response = salt['http.query'](url="https://api.github.com/orgs/tinkerbell/actions/runners/registration-token", method="POST", header="Accept: application/vnd.github.v3+json", username=pillar['github']['username'], password=pillar['github']['accessToken']).body | load_json %}

runner-configure:
  cmd.run:
    # TODO: do we need to add a labels config here, or is the arch label automatic?
    - name: ./config.sh --url https://github.com/tinkerbell --token {{ response['token'] }} --unattended --replace {% if grains['gha_runner_states'] is defined %}--labels {{ grains['gha_runner_states'] | join(', ') }}{% endif %}
    - cwd: /opt/actions-runner
    - runas: github
    - require:
      - archive: runner-download

runner-install-service:
  cmd.run:
    - name: ./svc.sh install github
    - cwd: /opt/actions-runner
    - require:
      - cmd: runner-configure

runner-enable-service:
  service.running:
    - name: actions.runner.tinkerbell.{{ grains.nodename }}.service
    - enable: True
    - reload: True
