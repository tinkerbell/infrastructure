include:
  - teleport.install
    
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
