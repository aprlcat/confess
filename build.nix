{ buildGoModule }:
buildGoModule {
  src = ./.;

  name = "confess";
  vendorHash = "sha256-91b0Jn8YZPb/wbQO0R4OME6o1HyzCzhhSaDMtZ4JVU4=";

  ldflags = [
    "-s"
    "-w"
  ];

  installPhase = ''
    runHook preInstall

    mkdir -p $out/bin
    install -Dm755 "$GOPATH/bin/confess" -T $out/bin/confess

    cp -r static $out/static

    runHook postInstall
  '';

  meta = {
    description = "A simple confessional website";
    homepage = "https://forge.catnip.ee/batteredbunny/confess";
    mainProgram = "confess";
  };
}
