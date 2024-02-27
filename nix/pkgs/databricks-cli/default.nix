{ lib, stdenvNoCC, fetchurl, unzip, installShellFiles, autoPatchelfHook, openssl }:

let
  pname = "databricks-cli";
  version = "0.214.1";

  url = "https://github.com/databricks/cli/releases/download/v${version}";

  sources = {
    x86_64-linux = fetchurl {
      url = "${url}/databricks_cli_${version}_linux_amd64.zip";
      sha256 = "sha256-Wu5nFDRlBIWE9SROlaNh0jpGYTyMdH7VwbGYQzL7LCo=";
    };

    aarch64-linux = fetchurl {
     url = "${url}/databricks_cli_${version}_linux_arm64.zip";
     sha256 = "sha256-yYRVLP8/RR8onV7mPA/Mkbxrke30ZwC9b/jaKCmkcUc=";
    };

    x86_64-darwin = fetchurl {
      url = "${url}/databricks_cli_${version}_darwin_amd64.zip";
      sha256 = "sha256-npw1AOr3hGn2uF0K9DxRX6N08aVFwCdk8SnXGoLQLJ0=";
    };
    aarch64-darwin = fetchurl {
      url = "${url}/databricks_cli_${version}_darwin_arm64.zip";
      sha256 = "sha256-N79jB+BXopQ+roBtuAEH1EzVjId9KZdcBTpyIyaeevA=";
    };
  };

  src = sources.${stdenvNoCC.hostPlatform.system} or (throw "Unsupported system: ${stdenvNoCC.hostPlatform.system}");
in
stdenvNoCC.mkDerivation rec {
  inherit pname version src;

  name = pname;

  strictDeps = true;
  nativeBuildInputs = [ unzip installShellFiles ] ++ lib.optionals stdenvNoCC.isLinux [ autoPatchelfHook ];
  buildInputs = [ openssl ];

  dontConfigure = true;
  dontBuild = true;

  unpackPhase = ''
    unzip $src
  '';

  installPhase = ''
    runHook preInstall

    install -Dm 755 ./databricks $out/bin/databricks

    runHook postInstall
  '';

  meta = with lib; {
    description = "A CLI client for Databricks";
    homepage = "https://github.com/databricks/cli";
    changelog = "https://github.com/databricks/cli/releases/tag/v${version}";
    license = "DB license";
    maintainers = [ "ahappypie" ];
    mainProgram = "databricks";
  };
}
