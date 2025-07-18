# vnc-summarizer

🌍 *[English](README.md) ∙ [Português](README_pt.md)*

`vnc-summarizer` é o serviço responsável pelo software que extrai os dados e sumariza as proposições disponibilizadas na
plataforma [Você na Câmara (VNC)](#você-na-câmara). Neste repositório você encontrará o código-fonte do software
principal da VNC, que utiliza tecnologias como Go, PostgreSQL e Amazon Web Services (AWS). Além disso, está disponível o
container Docker responsável por executar este código, permitindo que você execute o projeto de forma simples e rápida.

## Como Executar

### Pré-requisitos

Para executar corretamente o `vnc-summarizer` você precisará ter os containers dos serviços
[`vnc-databases`](https://github.com/devlucassantos/vnc-databases) e
[`vnc-pdf-content-extractor-api`](https://github.com/devlucassantos/vnc-pdf-content-extractor-api) em execução, de modo
que o container desta aplicação tenha acesso aos serviços necessários para a consulta e manipulação dos dados.

Além disso, você precisará preencher também algumas variáveis do arquivo `.env`, localizado no diretório _config_
(`./src/config/.env`). Neste arquivo, você notará que algumas variáveis já estão preenchidas — isso ocorre porque se
referem a configurações padrão, que podem ser utilizadas caso você opte por não modificar nenhum dos containers
pré-configurados para rodar os repositórios que compõem a VNC. No entanto, sinta-se à vontade para alterar qualquer uma
dessas variáveis, caso deseje adaptar o projeto ao seu ambiente. Observe também que algumas destas variáveis não estão
preenchidas - isso ocorre porque seu uso está vinculado a contas específicas de cada usuário em plataformas externas ao
VNC e, portanto, seus valores devem ser fornecidos individualmente por quem deseja utilizar esses recursos. Essas
variáveis são:

* `AWS_REGION` → Região/Servidor da conta do usuário do IAM na AWS (Para uma explicação mais detalhada sobre as
  credenciais de acesso da AWS, acesse a [documentação oficial da AWS sobre o gerenciamento de chaves de acesso de
  usuários do IAM](https://docs.aws.amazon.com/pt_br/IAM/latest/UserGuide/id_credentials_access-keys.html))
* `AWS_ACCESS_KEY_ID` → ID de acesso do usuário do IAM na AWS
* `AWS_SECRET_ACCESS_KEY` → Chave secreta de acesso do usuário do IAM na AWS
* `AWS_S3_BUCKET` → Nome do bucket onde as imagens das proposições serão salvas no AWS S3
* `OPENAI_API_KEY` → Para o preenchimento desta variável deve-se [criar uma chave de API no ChatGPT](https://platform.openai.com/account/api-keys),
  serviço de IA atualmente utilizado pelo VNC

### Executando via Docker

Para executar o serviço, você precisará ter o [Docker](https://www.docker.com) instalado na sua máquina e executar o
seguinte comando no diretório raiz deste projeto:

````shell
docker compose up --build
````

### Documentação

Após a execução do projeto, o serviço iniciará a busca pelos dados legislativos e a sumarização das proposições, sendo
possível acompanhar todo esse processo por meio dos logs do próprio container.

## Você na Câmara

Você na Câmara (VNC) é uma plataforma de notícias desenvolvida para simplificar e tornar acessíveis às proposições
legislativas que tramitam na Câmara dos Deputados do Brasil. Por meio do uso de Inteligência Artificial, a plataforma
sintetiza o conteúdo desses documentos legislativos, transformando informações técnicas e complexas em resumos objetivos
e claros para a população em geral.

Este projeto integra o Trabalho de Conclusão de Curso dos desenvolvedores da plataforma e foi concebido com base
em arquiteturas como a hexagonal e a de microsserviços. A solução foi organizada em diversos repositórios, cada um com
responsabilidades específicas dentro do sistema:

* [`vnc-databases`](https://github.com/devlucassantos/vnc-databases): Responsável por gerenciar a infraestrutura de
  dados da plataforma. Principais tecnologias utilizadas: PostgreSQL, Redis, Liquibase e Docker.
* [`vnc-pdf-content-extractor-api`](https://github.com/devlucassantos/vnc-pdf-content-extractor-api): Responsável por
  realizar a extração de conteúdo dos PDFs utilizados pela plataforma. Principais tecnologias utilizadas: Python,
  FastAPI e Docker.
* [`vnc-domains`](https://github.com/devlucassantos/vnc-domains): Responsável por centralizar os domínios e regras de
  negócio da plataforma. Principal tecnologia utilizada: Go.
* [`vnc-summarizer`](https://github.com/devlucassantos/vnc-summarizer): Responsável pelo software que extrai os dados e
  sumariza as proposições disponibilizadas na plataforma. Principais tecnologias utilizadas: Go, PostgreSQL, Amazon Web
  Services (AWS) e Docker.
* [`vnc-api`](https://github.com/devlucassantos/vnc-api): Responsável por disponibilizar os dados para o frontend da
  plataforma. Principais tecnologias utilizadas: Go, Echo, PostgreSQL, Redis e Docker.
* [`vnc-web-ui`](https://github.com/devlucassantos/vnc-web-ui): Responsável por fornecer a interface web da plataforma.
  Principais tecnologias utilizadas: TypeScript, SCSS, React, Vite e Docker.
