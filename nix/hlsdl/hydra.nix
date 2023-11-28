{ inputs, cell }:

let
  pkgs = inputs.nixpkgs;
  rev = inputs.self.rev or "unknown";
in
{
  hlsdl = cell.apps.default;
}
