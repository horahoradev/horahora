FROM node:14-alpine

RUN apk add --update npm

WORKDIR /webapp

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["sh", "/entrypoint.sh"]
