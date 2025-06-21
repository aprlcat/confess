{ buildGoModule }:
buildGoModule {
  src = ./.;

  name = "confess";
  vendorHash = "sha256-JeY3oe9XZfLYqyHPSZtvdH/Y05TzC0hNsuooYx3QSUs=";

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
