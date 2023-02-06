# Choose whatever you want, version >= 1.16
FROM golang:1.19-alpine

WORKDIR /app


COPY . /app
COPY .air.toml .
RUN go mod download 
RUN go install github.com/cosmtrek/air@latest
CMD ["air", "-c", "/app/.air.toml"]

