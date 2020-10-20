FROM golang:1.14

WORKDIR /go/src/github.com/giulianos/mutants
COPY . .

RUN go install -v ./cmd/analyzer
EXPOSE 80

CMD ["analyzer"]
