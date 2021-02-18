docker-deps:
  pkg.installed:
    - pkgs:
      - apt-transport-https
      - ca-certificates
      - curl
      - gnupg-agent
      - software-properties-common

docker-repo:
  pkgrepo.managed:
    - humanname: Docker Engine
    - name: deb [arch={{ grains.osarch }}] https://download.docker.com/linux/ubuntu {{ grains.oscodename }} stable
    - file: /etc/apt/sources.list.d/docker-engine.list
    - gpgcheck: 1
    - key_url: https://download.docker.com/linux/ubuntu/gpg

docker-install:
  pkg.installed:
    - pkgs:
      - docker-ce
      - docker-ce-cli
      - containerd.io

docker-service:
  service.running:
    - name: docker
    - enable: True
    - reload: True

containerd-service:
  service.running:
    - name: containerd
    - enable: True
    - reload: True
