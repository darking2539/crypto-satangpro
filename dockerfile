FROM golang:1.19-buster AS builder
LABEL stage=builder
COPY . /app
WORKDIR /app
RUN go mod tidy && go mod download && go build .

FROM ubuntu:18.04
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /app
COPY --from=builder /app/crypto-satangpro /app
ENTRYPOINT ["./crypto-satangpro"]