inputs: { pkgs
        , config ? pkgs.config
        , lib ? pkgs.lib
        , self
        , system
        , ...
        }:
let
  cfg = config.services.confess-web;
in
{
  options.services.confess-web = {
    enable = lib.mkEnableOption "confess-web";

    package = lib.mkOption {
      description = "package to use";
      default = inputs.self.packages.${system}.default;
    };

    dataDir = lib.mkOption {
      type = lib.types.str;
      default = "/var/lib/confess-web";
      description = "The directory where confess stores its database.";
    };

    user = lib.mkOption {
      type = lib.types.str;
      default = "confess-web";
      description = "User account under which confess runs.";
    };

    group = lib.mkOption {
      type = lib.types.str;
      default = "confess-web";
      description = "Group under which confess runs.";
    };

    port = lib.mkOption {
      type = lib.types.int;
      description = "port to run http api on";
    };

    ntfy = lib.mkOption {
      type = lib.types.nullOr lib.types.str;
      default = null;
      description = "ntfy url to use";
    };
  };

  config = lib.mkIf cfg.enable {
    systemd.tmpfiles.rules = [
      "d '${cfg.dataDir}' 0700 ${cfg.user} ${cfg.group} - -"
    ];

    users.users = lib.mkIf (cfg.user == "confess-web") {
      confess-web = {
        group = cfg.group;
        home = cfg.dataDir;
        createHome = true;
        isSystemUser = true;
      };
    };

    users.groups.${cfg.group} = { };

    systemd.services.confess-web = {
      enable = true;
      serviceConfig = {
        User = cfg.user;
        Group = cfg.group;
        ProtectSystem = "full";
        ProtectHome = "yes";
        DeviceAllow = [ "" ];
        LockPersonality = true;
        MemoryDenyWriteExecute = true;
        PrivateDevices = true;
        ProtectClock = true;
        ProtectControlGroups = true;
        ProtectHostname = true;
        ProtectKernelLogs = true;
        ProtectKernelModules = true;
        ProtectKernelTunables = true;
        ProtectProc = "invisible";
        RestrictNamespaces = true;
        RestrictRealtime = true;
        RestrictSUIDSGID = true;
        SystemCallArchitectures = "native";
        PrivateUsers = true;
        StateDirectory = cfg.dataDir;
        ExecStart = "${lib.getExe cfg.package} --port=${toString cfg.port} --static=${cfg.package}/static --database=${cfg.dataDir}/confess.db ${lib.optionalString (!isNull cfg.ntfy) "--ntfy=${cfg.ntfy}"}";
        Restart = "always";
      };
      wantedBy = [ "default.target" ];
    };
  };
}
