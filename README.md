# confess
A simple confessional website written in golang

```
nix run git+https://forge.catnip.ee/batteredbunny/confess
```


## Running as service on nixos
```nix
# flake.nix
inputs = {
    confess.url = "git+https://forge.catnip.ee/batteredbunny/confess";
};
```

```nix
# configuration.nix
imports = [
    inputs.confess.nixosModules.default
];

services.confess-web = {
    enable = true;
    port = 8080; # port to run http api on

    # Optional parameters
    reverseProxy = false; # enable if running behind reverse proxy
    trustedProxy = "127.0.0.1"; # Reverse proxy to trust
    package = inputs.lastfm-status.packages.${builtins.currentSystem}.default;
    ntfyUrl = ""; # ntfy url to use
    user = "confess-web"; # User account under which confess runs.
    group = "confess-web"; # Group under which confess runs.
    environmentFile = "/etc/secrets/confess.env"; # Useful for storing api keys like: NTFY_URL
};
```