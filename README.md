# Trabalho Semestral - Algoritmos e Estruturas de Dados III

Trabalho desenvolvido ao longo da disciplina de AED's III do curso de Bacharelado em Ciência da Computação na PUC Minas durante o primeiro semestre de 2023.

## Integrantes
* Bernardo Marques Fernandes
* Marcos Antônio Lommez Cândido Ribeiro

## Descrição do Projeto
O objetivo do projeto é criar um banco de dados em arquivo binário, utilizando técnicas de CRUD, ordenação, indexação, compactação, casamento de padroes e criptografia.
Além disso, o projeto foi feito com o modelo de uma API para comunicação entre backend e frontend.

## Funcionalidades
O banco de dados suporta as seguintes funcionalidades:

* Importação do csv
* CRUD
* Ordenação externa
* Indexação
* Pesquisa inclusiva
* Casamento de padroes
* Compactação
* Criptografia
  
## Exemplos de telas do sistema:

#### Tela Principal
![Tela principal](/Outros/Tela_Inicial.png)
#### Tela de visualização individual e crud
![PokeCard](/Outros/Dados.png)
#### Tela de pesquisa
![Resposta do servidor](/Outros/Pesquisa.png)
#### Menus de controles disponiveis
![Resposta do servidor](/Outros/Indexacao.png)
#### Tempo de execucao dos algoritmos
![Resposta do servidor](/Outros/Tempo.png)
#### Tela de respostas do servidor
![Resposta do servidor](/Outros/Carregamento.png)

## Tecnologias utilizadas
O projeto foi desenvolvido em linguagem Go e utiliza apenas a biblioteca padrão para todas as operações.
Frontend feito em javascript puro e sem framework.

## Algoritmos implementados
* Intercalação comum
* Intercalação com bloco de tamanho variavel
* Intercalação com substituição (Heap minimo)
* Hash dinamico
* Arvore B
* Arvore B+
* Indice Invertido
* KMP
* RabinKarp
* Huffman
* LZW
* AES CBC 128, 196, 256
* Trivium

## Videos de apresentações
* [TP1 - CRUD/Ordenacao](https://youtu.be/t9WriRSQGYM)
* [TP2 - Indexacao](https://youtu.be/VZeUh_TTPIE)
* [TP3 - Casamento de Padrões](https://youtu.be/FU3NHFADTt4)
* [TP4 - Compactação](https://youtu.be/JAGKt8K1VgQ)
* [TP5 - Criptografia](https://youtu.be/G9wO67tj6pA)

## Como utilizar
Para utilizar o banco de dados, basta baixar ou clonar o repositório e compilar o código-fonte na pasta **backend** com o comando:

`go run main.go`

Em seguida, é possível executar o programa abrindo o arquivo **index.html** na pasta **frontend**

## Contribuição
Esse projeto foi desenvolvido como trabalho de faculdade e portifolio, portanto não está aberto a contribuições externas. No entanto, sinta-se livre para utilizar o código como referência ou para fins educacionais.
