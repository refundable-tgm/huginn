FROM golang:1.14

WORKDIR /go/src/huginn
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

CMD ["huginn"]
