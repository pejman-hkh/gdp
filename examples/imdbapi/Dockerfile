FROM golang:1.22
#ENV http_proxy socks5://127.0.0.1:1088
#ENV https_proxy socks5://127.0.0.1:1088
WORKDIR /usr/src/app
COPY go.mod go.sum ./
COPY . .
RUN go get && go build -v -o /usr/local/bin/app ./...

EXPOSE 8090
CMD ["app"]