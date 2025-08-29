FROM golang:1.22 as builder
WORKDIR /src
COPY . .
RUN go build -o /gitwo ./cmd/gitwo

FROM alpine:3.19
COPY --from=builder /gitwo /usr/local/bin/gitwo
ENTRYPOINT ["gitwo"]

