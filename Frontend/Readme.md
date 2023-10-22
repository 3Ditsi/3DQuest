# 3DQuestFrontend

# Configuración y uso

## Instalación y configuración inicial

### Usuarios de Windows

Instalaremos Node y npm mediante el siguiente instalador: https://nodejs.org/en/download/ ### Usuarios de Linux Abre una terminal. Instala nodejs and npm: ```bash sudo apt install nodejs sudo apt install npm ``` Instala un gestor de versiones de node: ```bash sudo npm install -g n ``` Escoge node v16: ```bash sudo n 16
```

Ve a la carpeta 3DQuestFrontend/3DQuestFrontend.
Asegurate de instalar las dependencias:

```bash
yarn install
```

## Servidor de desarrollo

Lanza el servidor de desarrollo en http://localhost:3000

```bash
yarn run dev
```

o

```bash
yarn dev
```

## Producción

Compila la aplicación para producción:

```bash
yarn run build
```

---


# English
# Configuration and use

### Windows users

We will install Node and npm with the following installer: https://nodejs.org/en/download/ 

## Setup
Install nodejs and npm:
```bash
sudo apt install nodejs
sudo apt install npm
```
Install manager for node versions:

```bash
sudo npm install -g n
```

Choose node v16:
```bash
sudo n 16
```

Open a terminal and go to 3DQuestFrontend/3DQuestFrontend folder.
Make sure to install the dependencies:

```bash
yarn install

```

## Development Server

Start the development server on http://localhost:3000

```bash
yarn run dev
```

or

```bash
yarn dev
```

## Production

Build the application for production:

```bash
yarn run build
```
