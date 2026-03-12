FROM golang:1.25-alpine AS build
WORKDIR /www/app
COPY go.mod go.sum ./
RUN ["go", "mod", "download"]
COPY . .
RUN ["go", "build", "-o", "/www/app/main", "./cmd/api"]
FROM alpine
RUN apk add --no-cache bash
WORKDIR /www/app
COPY --from=build /www/app/main .
CMD ["./main"]
