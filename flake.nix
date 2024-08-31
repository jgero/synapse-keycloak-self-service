{
  description = "Matrix self service open id account creator";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, treefmt-nix }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
    in
    {
      formatter.${system} = treefmtEval.config.build.wrapper;
      checks.${system}.formatter = treefmtEval.config.build.check self;
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
        ];
        shellHook = ''
          export KEYCLOAK_HOST="auth.sfz-aalen.space"
          export KEYCLOAK_REALM="master"
          export KEYCLOAK_CLIENT_ID="matrix-self-service"
          export SYNAPSE_HOST="matrix.aalen.space"
          export SYNAPSE_ACCESS_TOKEN="<token>"
          export SYNAPSE_DOMAIN="aalen.space"
        '';
      };
      packages.${system} = rec {
        bin = pkgs.buildGoModule {
          pname = "synapse-keycloak-self-service";
          version = "1.0";
          # vendorHash = nixpkgs.lib.fakeHash;
          vendorHash = "sha256-KqkzYjhsgmD+BKZa/G28GRcMI63xcril+LCMpiYZMZE=";
          src = ./.;
        };
        default = pkgs.dockerTools.buildLayeredImage {
          name = "ghcr.io/jgero/${bin.pname}";
          tag = "latest";
          contents = with pkgs; [ cacert bin ];
          maxLayers = 10;
          config = {
            Cmd = [ "${bin}/bin/${bin.pname}" ];
            Env = with pkgs; [
              "SSL_CERT_FILE=${cacert}/etc/ssl/certs/ca-bundle.crt"
            ];
          };
        };
      };
    };
}
