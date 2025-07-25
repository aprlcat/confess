{ buildGoModule }:
buildGoModule {
  src = ./.;

  name = "confess";
  vendorHash = "sha256-PhrJpzMI8oVyF6qnkuMjNswscthkHyb+4Iinuzbv2G8=";

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
    homepage = "https://github.com/BatteredBunny/confess";
    mainProgram = "confess";
  };
}
