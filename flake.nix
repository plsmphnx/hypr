{
  outputs = { nixpkgs, ... }: let
    systems = fn: nixpkgs.lib.mapAttrs (_: fn) nixpkgs.legacyPackages;
    build = pkgs: name: let 
      bin = "hypr${name}";
    in pkgs.stdenv.mkDerivation {
      name = bin;
      src = ./.;
      nativeBuildInputs = [ pkgs.go ];
      buildPhase = "HOME=. go build -o ${bin} ./${name}";
      installPhase = "mkdir -p $out/bin && mv ${bin} $out/bin";
    };
  in {
    packages = systems (pkgs: {
      jump = build pkgs "jump";
      keys = build pkgs "keys";
    });
  };
}
