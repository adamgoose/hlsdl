{ inputs, cell }: {

  default = { lib, config, ... }:
    let
      cfg = config.services.hlsdl;
      l = lib // builtins;
    in
    {
      options.services.hlsdl = {
        enable = l.mkEnableOption (l.mdDoc "HLSDL Server");

        user = l.mkOption {
          type = l.types.str;
          description = "User to run the HLSDL Server as";
        };

        redisPort = l.mkOption {
          type = l.types.int;
          default = 6379;
          description = "Redis port";
        };

        redisDb = l.mkOption {
          type = l.types.int;
          default = 0;
          description = "Redis database";
        };
      };

      config = l.mkIf cfg.enable {

        systemd.packages = [
          cell.apps.default
        ];

        environment.systemPackages = [
          cell.apps.default
        ];

        # Run a Redis server
        services.redis.servers.hlsdl = {
          enable = true;
          port = cfg.redisPort;
        };

        # Run the HLSDL server
        systemd.services.hlsdl = {
          description = "HLSDL Server";
          after = [ "network.target" ];
          wantedBy = [ "multi-user.target" ];

          environment = {
            HLSDL_REDIS_ADDR = "localhost:${l.toString cfg.redisPort}";
            HLSDL_REDIS_DB = l.toString cfg.redisDb;
          };

          serviceConfig = {
            Restart = "on-failure";
            SuccessExitStatus = "3 4";
            RestartForceExitStatus = "3 4";
            User = cfg.user;
            ExecStart = "${cell.apps.default}/bin/hlsdl server";
            WorkingDirectory = "/tmp";
          };
        };

      };
    };

}
