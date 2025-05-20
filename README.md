# PixelForging
<p align="center">
  <img src="tests/input/Logo-Pixel-Forging.png" alt="Logo do Projeto Pixel Forging Versão Inicial" width="200">
</p>
PixelForging é uma ferramenta de linha de comando (CLI) e servidor gRPC para processamento digital de imagens, com foco em extração de paletas de cores.

## Funcionalidades

- **Extração de paleta de cores** de imagens, com opções para configurar o formato da paleta.
- **Servidor gRPC** para integração com outros sistemas.
- CLI amigável para uso rápido e scripts.

## Instalação

Clone o repositório e instale as dependências Go:

```sh
git clone https://github.com/Joao-lucas-felix/PixelForging.git
cd PixelForging
go mod tidy
go build 
```
## Como Usar

Para usar a CLI do PixelForging você pode fazer as seguintes ações: 

- Extrair Paletas de Cores.
- Iniciar um servidor gRPC.
- Dar um Hello para verificar se o binário está bem formado.

### Hello.

Hello saúda o usuário, serve para verificar a integridade da CLI. Você pode passar a flag Name para a CLI lhe saudar por nome. 

```bash
./Pixelforging hello 
```

A CLI retornara:  Hello  World

```bash
./Pixelforging hello --name João
```

A CLI retornara: Hello João. 

### Extract Palette

Abre a imagem no caminho informado pela flag --input-image="[CAMINHO_DA_IMAGEM]", extrai a paleta de cores e salva no caminho especificado pela flag --output-image="[CAMINHO_DA_PALETA]".

Você pode configurar o formato da imagem da paleta usando os seguintes parâmetros opcionais:
--colors-per-row="[NÚMERO_DE_CORES_POR_LINHA]"
--width="[LARGURA_DO_BLOCO_DE_COR]"
--height="[ALTURA_DO_BLOCO_DE_COR]"
--colors-num="[NÚMERO_TOTAL_DE_CORES]"

Valores padrão:
--colors-per-row=3
--width=50
--height=50
--colors-num=6

```bash
#Extrair a paleta de cores de uma imagem (usando os valores padrão)
./PixelForging extract-palette 
	--input-image tests/input/image.png 
	--output-image tests/out/palette_default.png
	
#Extrair a paleta de cores definindo o número de cores por linha
./PixelForging extract-palette 
	--input-image tests/input/image.png 
	--output-image tests/out/palette_row4.png 
	--colors-per-row 4

#Extrair a paleta de cores definindo largura e altura dos blocos de cor
./PixelForging extract-palette 
	--input-image tests/input/image.png 
	--output-image tests/out/palette_customsize.png 
	--width 60 
	--height 60
#Extrair a paleta de cores definindo o número total de cores
./PixelForging extract-palette 
	--input-image tests/input/image.png 
	--output-image tests/out/palette_12colors.png 
	--colors-num 12
```

### Exemplo de extração de Paleta:

<p align="center">
  <img src="tests/input/Logo-Pixel-Forging.png" alt="Logo do Projeto Pixel Forging Versão Inicial" width="200">
  <img src="tests/out/output2.png" alt="Paleta de cores do Projeto Pixel Forging Versão Inicial" width="200">
</p>


### Serviço gRPC

O serviço gRPC serve para que você seja capaz de usar as funções do PixelForging através da rede usando o protocolo HTTP. Usando a capacidade de Streaming bidirecional do gRPC para otimizar o trafego das imagens de entrada e saída pela rede. 

O serviço PixelForging está definido sobre os seguinte .proto: 

```protobuf
syntax = "proto3";

package pixelforging_grpc;

option go_package = "./src/backend/pb/pixelforging-grpc";

service PixelForging {
    rpc ExtractPalette(stream ExtractPaletteInput) returns (stream ExtractPaletteOutput);
}

message ExtractPaletteInput {
    bytes fileBytes = 1;
    string fileName = 2;
    string fileType = 3;
    // The following fields are optional and can be set to 0 if not needed
    // The following fields configure shape of the palette
    int32 colorsPerRow = 4;
    int32 colorWidth = 5; 
    int32 colorHeight = 6;
    int32 colorNum = 7;
}

message ExtractPaletteOutput {
    bytes paletteBytes = 1;
    string fileName = 2;
    string fileType = 3; 
}
```

### Para iniciar o servidor gRPC na porta padrão usando o binario do PixelForging:

```bash
./PixelForging start-gRPC-server
```

### Iniciar o servidor gRPC em uma porta específica

```bash
./PixelForging start-gRPC-server --port 50051
```
