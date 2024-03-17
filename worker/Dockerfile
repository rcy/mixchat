FROM node:16-bullseye

RUN apt-get update && apt-get install -y ffmpeg bsdgames python-is-python3
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && chmod a+rx /usr/local/bin/yt-dlp

WORKDIR /app

COPY . .

RUN npm install

EXPOSE 3010
CMD ["node", "./worker.js"]
