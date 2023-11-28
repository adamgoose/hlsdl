{ inputs, cell }:
let
  inherit (inputs) devenv n2c cells;
  pkgs = cell.pkgs.default;
in
{
  default = devenv.lib.mkShell {
    inherit inputs pkgs;
    modules = [
      ({ pkgs, ... }: {

        languages.go.enable = true;

        packages = with pkgs; [
          gomod2nix
          ffmpeg_5-headless
        ];

        services.redis = {
          enable = true;
        };

        scripts.asynq.exec = ''
          ${cells.asynq.apps.asynq}/bin/asynq \
            --uri=localhost:6379 \
            --db=0 \
            $@
        '';

        pre-commit.hooks = {
          gomod2nix = {
            enable = true;
            entry = "${pkgs.gomod2nix}/bin/gomod2nix";
            files = "go.mod|go.sum";
            pass_filenames = false;
          };
        };

      })
    ];
  };
}
