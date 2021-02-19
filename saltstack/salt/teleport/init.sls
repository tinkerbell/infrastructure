teleport-install:
  pkg.installed:
    - sources:
      - teleport: https://get.gravitational.com/teleport_5.1.0_amd64.deb

teleport-config:
  file.managed:
    - name: /etc/teleport.yaml
    - source: salt://{{ slspath }}/files/teleport.yaml
    - template: jinja
    - owner: root
    - group: root
    - mode: 644

teleport-config-github:
  file.managed:
    - name: /etc/teleport.github.yaml
    - source: salt://{{ slspath }}/files/github.yaml
    - template: jinja
    - owner: root
    - group: root
    - mode: 0400

certbot:
  pkg.installed

include:
  - ./letsencrypt

teleport-service:
  service.running:
    - name: teleport
    - enable: True
    - reload: True

teleport-service-wait:
  cmd.run:
    - name: sleep 2
    - require:
      - service: teleport-service

teleport-github-auth:
  cmd.run:
    - name: tctl create -f /etc/teleport.github.yaml
    - require:
      - cmd: teleport-service-wait

# ssh-service:
#   service.dead:
#     - name: ssh
#     - enable: False
