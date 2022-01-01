FROM node:16-bullseye

RUN apt-get update && apt-get install -y ffmpeg bsdgames python-is-python3

WORKDIR /app

COPY . .

RUN npm install

EXPOSE 3010
CMD ["node", "./worker.js"]
