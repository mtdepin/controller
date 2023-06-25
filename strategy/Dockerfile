FROM ubuntu:20.04
RUN mkdir  -p /controller/strategy
WORKDIR    /controller/strategy

COPY  ./strategy /controller/strategy/
RUN chmod +x /controller/strategy/strategy
VOLUME /controller/strategy/conf
VOLUME /controller/strategy/logs


CMD ["./strategy"]
