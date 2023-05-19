# Trabalho Semestral - Algoritmos e Estruturas de Dados III

Trabalho desenvolvido ao longo da disciplina de AED's III do curso de Bacharelado em Ciência da Computação na PUC Minas.

## Integrantes
* Bernardo Marques Fernandes
* Marcos Antonio Lommez

## Descrição do Projeto
O objetivo do projeto é criar um banco de dados em arquivo binário, utilizando técnicas de CRUD, ordenação, indexação, compactação, casamento de padroes e criptografia.
Além disso, o projeto foi feito com o modelo de uma API para comunicação entre backend e frontend.

### Exemplos de telas do sistema:

#### Tela Principal
![Tela principal](/Outros/Tela_Inicial.png)
#### Tela de visualização individual e crud
![PokeCard](/Outros/Dados.png)
#### Tela de pesquisa
![Resposta do servidor](/Outros/Pesquisa.png)
#### Escolha de Indexação e tempo de pesquisa
![Resposta do servidor](/Outros/Indexação.png)
#### Tela de respostas do servidor
![Resposta do servidor](/Outros/Carregamento.png)

## Funcionalidades
O banco de dados suporta as seguintes funcionalidades:

* Importação para população inicial
* CRUD completo
* Ordenação externa com diferentes métodos
* Indexação em memoria secundaria com Hash, Arvore B e B+
* Pesquisa com filtragem, ordenação por relevancia e combinação de pesquisas

** A ser implementado:
* Compactação de arquivo
* Casamento de padroes
* Criptografia

## Tecnologias utilizadas
O projeto foi desenvolvido em linguagem Go e utiliza a biblioteca padrão do Go para manipulação de arquivos binários.
Frontend feito em javascript puro e sem framework.

## Como utilizar
Para utilizar o banco de dados, basta baixar ou clonar o repositório e compilar o código-fonte na pasta **backend** com o comando:

`go run main.go`

Em seguida, é possível executar o programa abrindo o arquivo **index.html** na pasta **frontend**

## Contribuição
Esse projeto foi desenvolvido como trabalho de faculdade e não está aberto a contribuições externas. No entanto, sinta-se livre para utilizar o código como referência ou para fins educacionais.
