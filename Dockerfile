FROM golang:1.16-alpine

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY podchaosmonkey.go ./
RUN go build -o /pod_chaos_monkey
CMD [ "/pod_chaos_monkey" ]
