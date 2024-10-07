<h1 align="center">
Planning Poker 🎲
</h1>

<p align="center">
  <strong>Planeje sua próxima Sprint votando nas tarefas!</strong>
</p>

## Importante

> **_Aviso Legal_**: Este projeto é um clone do [Planning Poker Online](https://planningpokeronline.com/) e foi criado apenas para fins de aprendizado. Não tem uso comercial e não deve ser utilizado para esses fins.

### Como é a aparência?

![Sessão de Votação](https://user-images.githubusercontent.com/62121154/166455320-ff1159d0-7369-4c7b-9f06-0e6d8c557988.png)
![Cartas Reveladas](https://user-images.githubusercontent.com/62121154/166455328-ec64189f-641e-4be0-815c-e704ab1b753e.png)

### Estrutura

| Base de Código   |                                   Descrição                                   |
| :--------------- | :---------------------------------------------------------------------------: |
| [api](api)       | API escrita em **Go**, usa **Redis** como banco de dados em memória e pub/sub |
| [webapp](webapp) |                   Frontend em **Next.js** com _Typescript_                    |

### Diagrama do Sistema

![diagrama](https://user-images.githubusercontent.com/62121154/171638575-a2ffe342-b2be-4405-8dbd-33e1caea4271.png)

### Desenvolvimento

Para rodar a API e o Redis, primeiro precisamos construir as dependências necessárias:

```sh
docker-compose build
```

Então, para executar:

```sh
docker-compose up
```

Para executar o frontend, você precisa instalar o Node.js e então executar os comandos abaixo no diretório `/webapp`:

```sh
# primeiro instale as bibliotecas
npm install

# então execute o webapp
npm run dev
```
