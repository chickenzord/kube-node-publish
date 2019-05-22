FROM alpine:3.7

COPY bin/kube-node-publish /bin/
CMD ["kube-node-publish"]
