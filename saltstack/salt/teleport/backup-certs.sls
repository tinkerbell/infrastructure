letsencrypt-create-backup:
  cmd.run:
    - name: tar -cvf /tmp/letsencrypt.tar ./letsencrypt
    - cwd: /etc

letsencrypt-upload-backup:
  module.run:
    - name: s3.put
    - bucket: {{ pillar['s3']['bucketName'] }}            
    - local_file: /tmp/letsencrypt.tar
    - path: {{ grains.nodename }}/letsencrypt.tar
