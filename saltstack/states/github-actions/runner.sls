github:
  user.present:
    - groups:
      - wheel

runner-mkdir:
  file.directory:
    - name: /opt/actions-runner
    - user: github
    - group: github
    - mode: 0755
    - makedirs: True
    - require:
      - user: github

runner-download:
  archive.extracted:
    - name: /opt/actions-runner
    - source: https://github.com/actions/runner/releases/download/v2.277.1/actions-runner-linux-x64-2.277.1.tar.gz
    - source_hash: sha256=02d710fc9e0008e641274bb7da7fde61f7c9aa1cbb541a2990d3450cc88f4e98
    - user: github
    - group: github
    - if_missing: /opt/actions-runner/config.sh
    - require:
      - user: github

{% set response = salt['http.query'](name="https://api.github.com/orgs/tinkerbell/actions/runners/registration-token", method="POST", header="Accept: application/vnd.github.v3+json", username=pillar['github']['username'], password=pillar['github']['accessToken']).body | load_json %}

runner-configure:
  cmd.run:
    - name: ./config.sh --url https://github.com/tinkerbell --token {{ response['token'] }}
    - cwd: /opt/actions-runner
    - user: github
    - require:
      - archive: runner-download

runner-install-service:
  cmd.run:
    - name: ./svc.sh install github
    - user: github
    - cwd: /opt/actions-runner
    - require:
      - cmd: runner-configure

runner-enable-service:
  service.running:
    - name: actions.runner.tinkerbell.{{ grains.nodename }}.service
    - enable: True
    - reload: True
