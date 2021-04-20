include:
  - docker
  
github:
  user.present:
    - groups:
      - docker

install:
  pkg.installed:
    - pkgs:
      - git-lfs

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
    {% if grains.osarch == "amd64" %}
    - source: https://github.com/actions/runner/releases/download/v2.277.1/actions-runner-linux-x64-2.277.1.tar.gz
    - source_hash: sha256=02d710fc9e0008e641274bb7da7fde61f7c9aa1cbb541a2990d3450cc88f4e98
    {% elif grains.osarch == "arm64" %}
    - source: https://github.com/actions/runner/releases/download/v2.277.1/actions-runner-linux-arm64-2.277.1.tar.gz
    - source_hash: sha256=a6aa6dd0ba217118ef2b4ea24e9e0a85b02b13c38052a5de0776d6ced3a79c64
    {% endif %}
    - user: github
    - group: github
    - if_missing: /opt/actions-runner/config.sh
    - require:
      - user: github
