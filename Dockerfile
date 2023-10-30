# use example
# docker buildx build  --platform=linux/amd64,linux/arm64 --push -t iansmith/ofcimg .
FROM ubuntu:mantic
COPY ofcimg ofcimg
RUN mkdir data

ENV CGO_ENABLED=1
EXPOSE 9000:9000
ENTRYPOINT ["/ofcimg"]
