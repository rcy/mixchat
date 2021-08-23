FROM node:16-alpine

# youtube-dl needs python, ffmpeg
RUN apk add --update --no-cache python3 ffmpeg && ln -sf python3 /usr/bin/python

WORKDIR /app

COPY . .

RUN npm install

EXPOSE 3010
CMD ["node", "./worker.js"]
