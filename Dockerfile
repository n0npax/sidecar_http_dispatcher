FROM golang:latest as builder
LABEL maintainer=marcin.niemria@gmail.com
LABEL author="Marcin Niemira <n0npax>"
RUN mkdir /builder
WORKDIR /builder
COPY cmd cmd
COPY go.mod .
COPY go.sum .
COPY pkg pkg
RUN CGO_ENABLED=0 go build -o app cmd/main.go

FROM alpine:latest
LABEL maintainer=marcin.niemria@gmail.com
LABEL author="Marcin Niemira <n0npax>"
ENV SIDECAR_PORT 5000
RUN mkdir -p /app/config
COPY --from=builder /builder/app /main
CMD ["/main"]
