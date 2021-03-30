{% set response = salt['http.query'](url="https://api.github.com/orgs/tinkerbell/actions/runners/registration-token", method="POST", header="Accept: application/vnd.github.v3+json", username=pillar['github']['username'], password=pillar['github']['accessToken']).body | load_json %}

runner-configure:
  cmd.run:
    - name: ./config.sh --url https://github.com/tinkerbell --token {{ response['token'] }} --unattended --replace {% if grains['gha_runner_states'] is defined %}--labels {{ grains['gha_runner_states'] | join(', ') }}{% endif %}
    - cwd: /opt/actions-runner
    - runas: github
    - unless:
      - fun: file.file_exists
        path: /etc/systemd/system/actions.runner.tinkerbell.{{ grains.nodename }}.service
    - require:
      - archive: runner-download

runner-install-service:
  cmd.run:
    - name: ./svc.sh install github
    - cwd: /opt/actions-runner
    - unless:
      - fun: file.file_exists
        path: /etc/systemd/system/actions.runner.tinkerbell.{{ grains.nodename }}.service
    - require:
      - cmd: runner-configure

runner-enable-service:
  service.running:
    - name: actions.runner.tinkerbell.{{ grains.nodename }}.service
    - enable: True
    - reload: True
