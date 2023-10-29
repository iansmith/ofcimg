# use example
# docker buildx build  --platform=linux/amd64,linux/arm64 --push -t iansmith/ofcimg .
FROM alpine:3.18
COPY ofcimg ofcimg
RUN apk update
RUN apk add git go delve
RUN mkdir data
ENTRYPOINT ["/ofcimg"]
