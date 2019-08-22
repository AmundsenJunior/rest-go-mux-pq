FROM golang:1.12
LABEL maintainer="scottedwardrussell@gmail.com"
WORKDIR /build
COPY . /build
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -x -installsuffix cgo -o rest-go-mux-pq .

FROM alpine:latest
WORKDIR /app
COPY --from=0 /build/rest-go-mux-pq .
EXPOSE 8000
ENTRYPOINT ["./rest-go-mux-pq"]

