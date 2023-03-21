FROM golang:1.20-alpine as build

WORKDIR /go/src/k8s-log-proxy
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /go/bin/k8s-log-proxy

FROM gcr.io/distroless/base
USER 10000:10000
COPY --from=build /go/bin/k8s-log-proxy /

EXPOSE 8080
CMD [ "/k8s-log-proxy" ]
