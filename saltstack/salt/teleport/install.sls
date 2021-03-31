# Doesn't support Arm64
# teleport-repo:
#   pkgrepo.managed:
#     - humanname: Teleport
#     - name: deb https://deb.releases.teleport.dev/ stable main
#     - key_url: https://deb.releases.teleport.dev/teleport-pubkey.asc
teleport-download:
  archive.extracted:
    - name: /opt
    {% if grains.osarch == "amd64" %}
    - source: https://get.gravitational.com/teleport-v6.0.1-linux-amd64-bin.tar.gz
    - source_hash: sha256=d8463472ba2cfe34c77357bf16c02c0f7a381a7610ede81224ee8d064f908177
    {% elif grains.osarch == "arm64" %}
    - source: https://get.gravitational.com/teleport-v6.0.1-linux-arm64-bin.tar.gz
    - source_hash: sha256=d3c98ddbffb219eaa4a89410ced10c7f6a481cc2e326d03a73e4eda3feac6c9c
    {% endif %}
    - if_missing: /opt/teleport

teleport-install:
  cmd.run:
    - name: /opt/teleport/install

teleport-install-service:
  file.managed:
    - name: /etc/systemd/system/teleport.service
    - source: /opt/teleport/examples/systemd/teleport.service

