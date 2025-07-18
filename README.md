# vnc-summarizer

🌍 *[English](README.md) ∙ [Português](README_pt.md)*

`vnc-summarizer` is the service responsible for extracting data and summarizing the legislative propositions for the
[Você na Câmara (VNC)](#você-na-câmara) platform. In this repository, you will find the source code of VNC’s core
software, which uses technologies such as Go, PostgreSQL, and Amazon Web Services (AWS). Additionally, the Docker
container responsible for running this code is available, allowing you to execute the project quickly and easily.

## How to run

### Prerequisites

To properly run `vnc-summarizer`, you will need to have the [`vnc-databases`](https://github.com/devlucassantos/vnc-databases)
and [`vnc-pdf-content-extractor-api`](https://github.com/devlucassantos/vnc-pdf-content-extractor-api) service
containers running, so that this application's container has access to the necessary services for querying and
manipulating data.

Additionally, you will also need to fill in some variables in the `.env` file, located in the _config_ directory
(`./src/config/.env`). In this file, you’ll notice that some variables are already filled in — this is because they
refer to default configurations, which can be used if you choose not to modify any of the pre-configured containers
used to run the repositories that make up VNC. However, feel free to change any of these variables if you wish to adapt
the project to your environment. Also note that some of these variables are not filled in — this is because their use is
tied to specific user accounts on platforms external to VNC, and therefore their values must be provided individually by
whoever intends to use these features. These variables are:

* `AWS_REGION` → Region/Server of the IAM user account in AWS (For a more detailed explanation about AWS access
  credentials, visit the [official AWS documentation on managing IAM user access
  keys](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html))
* `AWS_ACCESS_KEY_ID` → IAM user access ID in AWS
* `AWS_SECRET_ACCESS_KEY` → IAM user secret access key in AWS
* `AWS_S3_BUCKET` → Name of the bucket where the images of the proposals will be saved in AWS S3
* `OPENAI_API_KEY` → To fill in this variable an [API key must be created in ChatGPT](https://platform.openai.com/account/api-keys),
  the AI service currently used by VNC

### Running via Docker

To run the service, you will need to have [Docker](https://www.docker.com) installed on your machine and run the
following command in the root directory of this project:

````shell
docker compose up --build
````

### Documentation

After running the project, the service will start retrieving legislative data and summarizing the propositions. You can
monitor the entire process through the container logs.

## Você na Câmara

Você na Câmara (VNC) is a news platform developed to simplify and make accessible the legislative propositions being
processed in the Chamber of Deputies of Brazil. Through the use of Artificial Intelligence, the platform synthesizes the
content of these legislative documents, transforming technical and complex information into clear and objective
summaries for the general public.

This project is part of the Final Paper of the platform's developers and was conceived based on architectures such as
hexagonal and microservices. The solution was organized into several repositories, each with specific responsibilities
within the system:

* [`vnc-databases`](https://github.com/devlucassantos/vnc-databases): Responsible for managing the platform's data
  infrastructure. Main technologies used: PostgreSQL, Redis, Liquibase, and Docker.
* [`vnc-pdf-content-extractor-api`](https://github.com/devlucassantos/vnc-pdf-content-extractor-api): Responsible for
  extracting content from the PDFs used by the platform. Main technologies used: Python, FastAPI, and Docker.
* [`vnc-domains`](https://github.com/devlucassantos/vnc-domains): Responsible for centralizing the platform's domains
  and business logic. Main technology used: Go.
* [`vnc-summarizer`](https://github.com/devlucassantos/vnc-summarizer): Responsible for the software that extracts data
  and summarizes the propositions available on the platform. Main technologies used: Go, PostgreSQL,
  Amazon Web Services (AWS), and Docker.
* [`vnc-api`](https://github.com/devlucassantos/vnc-api): Responsible for providing data to the platform's frontend.
  Main technologies used: Go, Echo, PostgreSQL, Redis, and Docker.
* [`vnc-web-ui`](https://github.com/devlucassantos/vnc-web-ui): Responsible for providing the platform's web interface.
  Main technologies used: TypeScript, SCSS, React, Vite, and Docker.
