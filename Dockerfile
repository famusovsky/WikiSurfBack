FROM golang:1.22.1-alpine

WORKDIR /

ARG override_tables=false
ARG port=8888

ENV OVERRIDE=$override_tables
ENV PORT=$port

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o out ./cmd/web

EXPOSE ${port}

CMD [ "sh", "-c", "./out -override_tables=$OVERRIDE -addr=:$PORT" ]