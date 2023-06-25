FROM ubuntu:20.04
RUN mkdir  -p /controller/business
WORKDIR    /controller/business

COPY  ./business /controller/business/
RUN chmod +x /controller/business/business
VOLUME /controller/business/conf
VOLUME /controller/business/logs


CMD ["./business"]
