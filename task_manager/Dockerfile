FROM golang
WORKDIR /app
COPY docker_manager/ .
RUN go build -o task_manager .
CMD ["./task_manager"]