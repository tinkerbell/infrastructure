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

teleport-tls:
  cmd.run:
    - name: certbot certonly -m "support@tinkerbell.org" --standalone --agree-tos --preferred-challenges http -d {{ pillar.teleport.domain }} -n

teleport-service:
  service.running:
    - name: teleport
    - enable: True
    - reload: True

teleport-github-auth:
  cmd.run:
    - name: tctl create /etc/teleport.github.yaml

ssh-service:
  service.dead:
    - name: ssh
    - enable: False
