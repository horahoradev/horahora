FROM node:14-alpine

RUN apk add --update npm

WORKDIR /webapp

COPY . .

RUN cd /webapp && \ 
    npm install && \
    npm install -D tailwindcss@npm:@tailwindcss/postcss7-compat postcss@^7 autoprefixer@^9 && npm install @craco/craco

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["sh", "/entrypoint.sh"]
