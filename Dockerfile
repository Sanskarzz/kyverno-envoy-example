FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN export GONOSUMDB="github.com/kyverno/kyverno-json"

RUN go get "github.com/kyverno/kyverno-json"

RUN go mod download

RUN go mod tidy 

RUN go build -o main ./cmd/main.go

EXPOSE 9002

ENTRYPOINT ["./main"]