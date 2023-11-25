### Build application 
FROM golang:1.17-alpine as builder
LABEL maintainer="Tanmay"
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o redbook-be cmd/redbook/main.go

### Run application from scratch
FROM  alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /build
COPY . .
COPY --from=builder /build/redbook-be .
EXPOSE 8085
CMD ["./redbook-be"]