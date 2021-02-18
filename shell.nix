with import <nixpkgs> {};

stdenv.mkDerivation {
  name = "env";
  buildInputs = [
    ruby.devEnv
    git
    gnumake
    pulumi-bin
  ];

  shellHook = ''
    export DOCKER_BUILDKIT=0
  '';
}
