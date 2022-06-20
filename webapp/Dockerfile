FROM node:14-alpine

RUN apk add --update npm

WORKDIR /webapp

# https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#add-or-copy
COPY ["package.json", "package-lock.json", "/webapp"]

RUN npm install

COPY . .

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["sh", "/entrypoint.sh"]