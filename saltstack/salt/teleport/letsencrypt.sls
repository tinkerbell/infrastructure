awscli:
  pkg.installed

# Only sync S3 to Local if the local LetsEncrypt directory doesn't exist
teleport-s3-sync-down:
  cmd.run:
    - name: aws s3 sync s3://{{ pillar.aws.bucketName }}/{{ grains.nodename }}/letsencrypt /etc/letsencrypt
    - env:
      - AWS_ACCESS_KEY_ID: {{ pillar.aws.accessKeyID }}
      - AWS_SECRET_ACCESS_KEY: {{ pillar.aws.secretAccessKey }}
      - AWS_DEFAULT_REGION: {{ pillar.aws.bucketLocation }}
    - require:
      - pkg: awscli
    - unless:
      - ls /etc/letsencrypt

# If we don't have a cert for our domain on the host, we need to issue one
teleport-tls-new:
  cmd.run:
    - name: certbot certonly -m "support@tinkerbell.org" --standalone --agree-tos --preferred-challenges http -d {{ pillar.teleport.domain }} -n
    # S3 sync will only fail if the creds are wrong, not if the remote doesn't exist
    # so we can use a require here to ensure we only really request a cert if the sync
    # has atleast run.
    - require:
      - cmd: teleport-s3-sync-down
    - unless:
      - ls /etc/letsencrypt/live/{{ pillar.teleport.domain }}/cert.pem

teleport-s3-sync-up:
  cmd.run:
    - name: aws s3 sync /etc/letsencrypt s3://{{ pillar.aws.bucketName }}/letsencrypt/
    - env:
      - AWS_ACCESS_KEY_ID: {{ pillar.aws.accessKeyID }}
      - AWS_SECRET_ACCESS_KEY: {{ pillar.aws.secretAccessKey }}
      - AWS_DEFAULT_REGION: {{ pillar.aws.bucketLocation }}
    - watch:
      - cmd: teleport-tls-new
