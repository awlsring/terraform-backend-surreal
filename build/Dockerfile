FROM golang:1.19-alpine

WORKDIR /app

COPY ./ ./

COPY ./cmd/main.go .

RUN go mod download

RUN go build -o /server

EXPOSE 8033

CMD [ "/server" ]