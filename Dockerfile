FROM golang:buster
LABEL maintainer=marcin.niemria@gmail.com
LABEL author="Marcin Niemira <n0npax>"
WORKDIR /app
RUN mkdir /app/config
COPY . .
COPY _nginx.conf /etc/nginx/nginx.conf
ENV SIDECAR_PORT 5000
ENTRYPOINT go run main.go
