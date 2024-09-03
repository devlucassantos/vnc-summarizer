# vnc-summarizer

🌍 *[English](README.md) ∙ [Português](README_pt.md)*

`vnc-summarizer` é o repositório responsável pelo software que sumariza as proposições exibidas na plataforma
[Você na Câmara (VNC)](#você-na-câmara-vnc). Neste repositório você encontrará o código-fonte do software principal do
VNC e também o container responsável por executar este código, deste modo você poderá facilmente rodar o projeto.

## Como Executar

### Pré-requisitos

Para executar este código você precisará preencher alguns campos do arquivo `.env` presente no diretório _config_
(`./src/config/.env`). Neste arquivo você poderá observar que alguns campos já estão preenchidos, isto porque são
configurações padrões que poderão ser utilizadas caso você opte por não modificar nenhum dos containers pré-configurados
para rodar os repositórios que compõem o VNC, entretanto fique a vontade para modificar quaisquer uma dessas variáveis
de modo a fazer o projeto se adaptar ao seu ambiente. Observe também que algumas destas variáveis não estão preenchidas,
isto ocorre porque essas variáveis tem seu uso vinculado a conta de cada usuário em plataformas externas ao VNC e por
isso devem ter seus valores gerados por cada usuário que deseje utilizar as plataformas. Essas chaves são:
* `AWS_REGION` → Região/Servidor da conta do usuário do IAM na AWS (Para uma explicação mais detalhada sobre as credenciais
de acesso da AWS, acesse a [documentação oficial da AWS sobre o gerenciamento de chaves de acesso de usuários do 
IAM](https://docs.aws.amazon.com/pt_br/IAM/latest/UserGuide/id_credentials_access-keys.html)). 
* `AWS_ACCESS_KEY_ID` → ID de acesso do usuário do IAM na AWS.
* `AWS_SECRET_ACCESS_KEY` → Chave secreta de acesso do usuário do IAM na AWS.
* `AWS_S3_BUCKET` → Nome do bucket onde as imagens das proposições serão salvas no AWS S3.
* `OPENAI_API_KEY` → Para o preenchimento desta variável deve-se [criar uma chave de API no ChatGPT](https://platform.openai.com/account/api-keys), 
serviço de IA atualmente utilizado pelo VNC.
* `UNICLOUD_KEY` → Para o preenchimento desta variável deve-se [criar uma chave de API no UniCLOUD](https://cloud.unidoc.io/#/api-keys),
serviço de manipulação de PDFs atualmente utilizado pelo VNC.

> Observe que para executar corretamente o `vnc-summarizer` você precisará ter os [containers do `vnc-databases`](https://github.com/devlucassantos/vnc-databases)
em execução de modo que o container desta aplicação tenha acesso aos bancos de dados necessários para a consulta e
modificação dos dados.

### Executando via Docker

Para executar o programa você precisará ter o [Docker](https://www.docker.com) instalado na sua máquina e executar o
seguinte comando no diretório raiz deste projeto:

````shell
docker compose up
````

### Documentação

Após a execução do projeto, o programa iniciará a busca e sumarização das proposições, inserindo-as por fim no banco de
dados da plataforma.

## Você na Câmara (VNC)

Você na Câmara (VNC) é uma plataforma de notícias que busca simplificar as proposições que tramitam pela Câmara dos
Deputados do Brasil visando sintetizar as ideias destas proposições através do uso da Inteligência Artificial (IA)
de modo que estes documentos possam ter suas ideias expressas de maneira simples e objetiva para a população em geral.
