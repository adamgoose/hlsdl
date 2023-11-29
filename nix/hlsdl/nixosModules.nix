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

        listen = l.mkOption {
          type = l.types.str;
          default = ":8881";
          description = "Listen address for the HTTP Server";
        };

        out = l.mkOption {
          type = l.types.str;
          description = "Output directory";
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
            HLSDL_OUT = cfg.out;
            HLSDL_LISTEN = cfg.listen;
            HLSDL_REDIS_ADDR = "localhost:${l.toString cfg.redisPort}";
            HLSDL_REDIS_DB = l.toString cfg.redisDb;
          };

          serviceConfig = {
            Restart = "on-failure";
            SuccessExitStatus = "3 4";
            RestartForceExitStatus = "3 4";
            User = cfg.user;
            ExecStart = "${cell.apps.default}/bin/hlsdl server";
          };
        };

      };
    };

}
