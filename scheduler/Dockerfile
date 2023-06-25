FROM ubuntu:20.04
RUN mkdir  -p /controller/scheduler
WORKDIR    /controller/scheduler

COPY  ./scheduler /controller/scheduler/
RUN chmod +x /controller/scheduler/scheduler
VOLUME /controller/scheduler/conf
VOLUME /controller/scheduler/logs


CMD ["./scheduler"]