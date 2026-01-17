FROM golang:alpine

ENV PORT=8080

WORKDIR /app

COPY . .

RUN apk update && apk add --no-cache sqlite sqlite-libs gcc g++ musl-dev
RUN go mod download
RUN go mod tidy

RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/gofrs/uuid/v5
RUN go get golang.org/x/crypto@v0.28.0

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o echohubApp /app/cmd/api/main.go

EXPOSE 8080

CMD ["./echohubApp"]

# docker build -t echohub-community-app .
# docker images
# docker run -p 8080:8080 <img-name>
