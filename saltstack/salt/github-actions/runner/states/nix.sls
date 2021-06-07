/nix:
  file.directory:
    - name: /nix
    - user: github
    - group: github
    - dir_mode: 755

https://nixos.org/nix/install:
  cmd.script:
    - runas: github
    - unless:
      - fun: file.file_exists
        path: /home/github/.nix-profile/etc/profile.d/nix.sh
