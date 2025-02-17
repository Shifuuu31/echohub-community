FROM golang:alpine

ENV GOLANG_VERSION 1.24.0

ENV PORT=8080

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o forumApp /app/cmd/api/main.go

EXPOSE 8080

CMD ["./forumApp"]



