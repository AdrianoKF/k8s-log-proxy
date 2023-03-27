FROM golang:1.20 as build

WORKDIR /go/src/k8s-log-proxy
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /go/bin/k8s-log-proxy

FROM gcr.io/distroless/static-debian11
USER 10000:10000
COPY --from=build /go/bin/k8s-log-proxy /

EXPOSE 8080
CMD [ "/k8s-log-proxy" ]
