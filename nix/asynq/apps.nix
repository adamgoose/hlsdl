{ inputs, cell }:
let
  # The `inputs` attribute allows us to access all of our flake inputs.
  inherit (inputs) nixpkgs std;

  # This is a common idiom for combining lib with builtins.
  l = nixpkgs.lib // builtins;
in
{

  asynq = with nixpkgs; buildGoModule rec {
    pname = "asynq";
    version = "0.24.0";

    src = fetchFromGitHub
      {
        owner = "hibiken";
        repo = "asynq";
        rev = "v${version}";
        sha256 = "sha256-vFB4IULTeMUnnvRRJBiSsyK0zmzYcql4TcUfUVO5dOE=";
      } + "/tools";

    vendorSha256 = "sha256-x2wOfwHRdR8XpGkxnQbf9B6Th5jAhcaRieZ5yq8RUEM=";
    subPackages = [ "asynq" ];
  };

}
