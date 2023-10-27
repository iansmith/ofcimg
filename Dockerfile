FROM alpine:3.18
COPY ofcimg ofcimg
RUN mkdir data
CMD ["/ofcimg"]