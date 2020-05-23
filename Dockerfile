FROM golang:latest
LABEL maintainer=marcin.niemria@gmail.com
LABEL author="Marcin Niemira <n0npax>"
WORKDIR /app
COPY . .
RUN go build -o app ./...

FROM scratch
LABEL maintainer=marcin.niemria@gmail.com
LABEL author="Marcin Niemira <n0npax>"
ENV SIDECAR_PORT 5000
COPY --from=0 /app/app /app
ENTRYPOINT /app
