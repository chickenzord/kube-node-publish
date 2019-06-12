FROM alpine:3.7

RUN apk add --update --no-cache ca-certificates \
    && update-ca-certificates

COPY bin/kube-node-publish /bin/
CMD ["kube-node-publish"]
