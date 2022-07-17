FROM node:14-alpine

RUN npm install --global npm@8.12

WORKDIR /webapp

# https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#add-or-copy
COPY ["package.json", "package-lock.json", "/webapp/"]

RUN npm install

COPY . .

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["sh", "/entrypoint.sh"]
