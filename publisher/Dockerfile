FROM golang:1.21-alpine AS build

WORKDIR /app
COPY . .

ENV GOPROXY=https://goproxy.io,direct
ENV GOPRIVATE=git.mycompany.com,github.com/my/private

RUN go build -o app ./cmd/publisher

FROM alpine:3.18

WORKDIR /root/
COPY --from=build /app/app .

EXPOSE 8080
CMD ["./app"]
