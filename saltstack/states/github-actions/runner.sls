runner-mkdir:
  file.directory:
    - name: /opt/actions-runner
    - user: root
    - group: root
    - mode: 0755
    - makedirs: True

runner-download:
  archive.extracted:
    - name: /opt/actions-runner
    - source: https://github.com/actions/runner/releases/download/v2.277.1/actions-runner-linux-x64-2.277.1.tar.gz
    - user: root
    - group: root
    - if_missing: /opt/actions-runner/config.sh

runner-configure:
  cmd.run:
    - name: ./config.sh --url https://github.com/tinkerbell/infrastructure --token {{ pillar['github']['actions']['token] }}
    - cwd: /opt/actions-runner

runner-run:
  cmd.run:
    - name: ./run.sh
    - cwd: /opt/actions-runner
