FROM ubuntu:jammy
RUN apt-get update
RUN apt-get install -y build-essential golang
RUN apt install -y docker.io 

ENV CGO_ENABLED=1
RUN mkdir data
