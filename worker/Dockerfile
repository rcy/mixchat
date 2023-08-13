FROM node:16-bullseye

RUN apt-get update && apt-get install -y ffmpeg bsdgames python-is-python3 pip
RUN pip install yt-dlp

WORKDIR /app

COPY . .

RUN npm install

EXPOSE 3010
CMD ["node", "./worker.js"]
