FROM golang:1.22-bullseye AS builder
ARG rev=dev
WORKDIR /work
COPY . .
RUN go mod download
RUN go build -o main cmd/api/main.go
RUN make -C vbin clean yt-dlp

FROM debian:bullseye
ARG cookies_file
RUN apt-get update && apt-get install -y ca-certificates python3 ffmpeg
WORKDIR /work
COPY --from=builder /work/main  .
COPY --from=builder /work/vbin/ ./vbin/
COPY $cookies_file ./vbin/
EXPOSE 8080
CMD ["./main"]
