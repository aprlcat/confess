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
    port = 8080;

    # Optional parameters
    package = inputs.lastfm-status.packages.${builtins.currentSystem}.default;
    ntfy = "ntfy.sh/asdsad";
};
```