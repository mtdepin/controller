FROM ubuntu:20.04
RUN mkdir  -p /controller/task_tracker
WORKDIR    /controller/task_tracker

COPY  ./task_tracker /controller/task_tracker/
RUN chmod +x /controller/task_tracker/task_tracker
VOLUME /controller/task_tracker/conf
VOLUME /controller/task_tracker/logs


CMD ["./task_tracker"]