{
  outputs = { nixpkgs, ... }: let
    systems = fn: nixpkgs.lib.mapAttrs (_: fn) nixpkgs.legacyPackages;
    build = pkgs: name: let 
      bin = "hypr${name}";
    in pkgs.stdenv.mkDerivation {
      name = bin;
      src = ./.;
      nativeBuildInputs = [ pkgs.nim ];
      buildPhase = "HOME=. nim c -d:release -d:strip -o:${bin} ./${name}/main.nim";
      installPhase = "mkdir -p $out/bin && mv ${bin} $out/bin";
    };
  in {
    packages = systems (pkgs: {
      jump = build pkgs "jump";
      mods = build pkgs "mods";
    });
  };
}
