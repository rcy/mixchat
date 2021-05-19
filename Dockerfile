FROM node:16-alpine

# youtube-dl needs python
RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY ./server.js ./

EXPOSE 3010
CMD ["node", "./server.js"]
