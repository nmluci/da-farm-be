FROM golang:1.22.5 as build
WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/

RUN go mod download
RUN go mod tidy

COPY . /app/

RUN CGO_ENABLED=0 go build -o /app/main /app/cmd/server

# Deploy

FROM alpine:3.16.0
WORKDIR /app

EXPOSE 7780

RUN apk update
RUN apk add --no-cache tzdata
ENV cp /usr/share/zoneinfo/Asia/Makassar /etc/localtime
RUN echo "Asia/Makassar" > /etc/timezone

# COPY --from=build /app/config /app/config
COPY --from=build /app/main /app/main

CMD ["/app/main"]