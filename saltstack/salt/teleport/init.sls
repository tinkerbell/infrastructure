include:
  - teleport.install
  
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

awscli:
  pkg.installed

teleport-s3-sync-down:
  archive.extracted:
    - name: /etc
    - source: s3://{{ pillar['s3']['bucketName'] }}/{{ grains.nodename }}/letsencrypt.tar
    - skip_verify: True
    - user: root
    - group: root
    - if_missing: /etc/letsencrypt

# If we don't have a cert for our domain on the host, we need to issue one
teleport-tls-new:
  cmd.run:
    - name: certbot certonly -m "support@tinkerbell.org" --standalone --agree-tos --preferred-challenges http -d {{ pillar.teleport.domain }} -n
    - unless:
      - ls /etc/letsencrypt/live/{{ pillar.teleport.domain }}/cert.pem

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

letsencrypt-backup:
  schedule.present:
    - function: state.sls
    - job_args:
      - teleport.backup-certs
    - cron: '0 3 * * *'
