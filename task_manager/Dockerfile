FROM golang:1.14

WORKDIR /app

COPY src/ .

RUN go build -o task_manager .

CMD ["./task_manager"]

EXPOSE 3333

ENV ELASTICSEARCH_HOST 127.0.0.1

# docker build -t task_manager .

# docker run -p 3333:3333 --name task_manager task_manager
