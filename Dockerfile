FROM golang:1.18

ENV CSV_UPLOAD_PORT="9100"
ENV PERSIST_URL="http://localhost:9001"

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . /usr/src/app
RUN go build -v -o /usr/local/bin/app

EXPOSE 9100

CMD ["app"]
