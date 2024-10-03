FROM golang:1.23-alpine AS build

WORKDIR /order

COPY go.mod go.sum ./

RUN go mod download


COPY . .

RUN go build -o main .


FROM alpine:latest

WORKDIR /order

COPY --from=build /order/main .

EXPOSE 5002

CMD ["./cmd/main"]
