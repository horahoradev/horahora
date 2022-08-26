FROM node:18.7-buster

RUN apk update && apk add python3 && npm install --global npm@8.17

WORKDIR /webapp

# https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#add-or-copy
COPY ["package.json", "package-lock.json", "/webapp/"]

RUN npm install

COPY . .

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["sh", "/entrypoint.sh"]
