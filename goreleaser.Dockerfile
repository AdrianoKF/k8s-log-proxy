FROM gcr.io/distroless/static-debian11
USER 10000:10000
COPY k8s-log-proxy /

EXPOSE 8080
CMD [ "/k8s-log-proxy" ]
