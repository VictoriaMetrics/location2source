FROM golang:1.25.5-alpine3.23 AS builder

WORKDIR /build
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o app ./...

FROM public.ecr.aws/docker/library/alpine:3.20

COPY --from=builder /build/app /app

CMD ["/app"]