FROM node:14-alpine

RUN apk add --update npm

WORKDIR /webapp

# https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#add-or-copy
COPY ["package.json", "yarn.lock", "/webapp"]

RUN cd /webapp && \ 
    npm install && \
    npm install -D tailwindcss@npm:@tailwindcss/postcss7-compat postcss@^7 autoprefixer@^9 && \
    npm install @craco/craco

COPY . .

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["sh", "/entrypoint.sh"]