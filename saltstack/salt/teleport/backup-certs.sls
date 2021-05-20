letsencrypt-create-backup:
    cmd.run:
        - name: tar -cvf /tmp/letsencrypt.tar /etc/letsencrypt

letsencrypt-s3-sync-up:
  file.managed:
    - name: s3://{{ pillar['s3.bucketName'] }}/{{ grains.nodename }}/letsencrypt.tar
    - source: /tmp/letsencrypt.tar
