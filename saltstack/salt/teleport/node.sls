teleport-install:
  pkg.installed:
    - sources:
      - teleport: https://get.gravitational.com/teleport_5.1.0_amd64.deb

teleport-config:
  file.managed:
    - name: /etc/teleport.yaml
    - source: salt://{{ slspath }}/files/node.yaml
    - template: jinja
    - owner: root
    - group: root
    - mode: 644

teleport-service:
  service.running:
    - name: teleport
    - enable: True
    - reload: True
