FROM golang:1.17-buster as build

WORKDIR /go/src/k8s-log-proxy
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /go/bin/k8s-log-proxy

FROM gcr.io/distroless/base
COPY --from=build /go/bin/k8s-log-proxy /

EXPOSE 8080
CMD [ "/k8s-log-proxy" ]
