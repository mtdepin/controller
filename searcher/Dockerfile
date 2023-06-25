FROM ubuntu:20.04
RUN mkdir  -p /controller/searcher
WORKDIR    /controller/searcher

COPY  ./searcher /controller/searcher/
RUN chmod +x /controller/searcher/searcher
VOLUME /controller/searcher/conf
VOLUME /controller/searcher/logs


CMD ["./searcher"]