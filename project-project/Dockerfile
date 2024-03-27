FROM alpine
WORKDIR /Initial
RUN mkdir -p /msproject/project-project/logs/debug
RUN mkdir -p /msproject/project-project/logs/info
RUN mkdir -p /msproject/project-project/logs/error
COPY ./target/project-project .
COPY ./config/bootstrap.yaml .
RUN  mkdir config && mv bootstrap.yaml config/bootstrap.yaml
EXPOSE 8081 8882
ENTRYPOINT ["./project-project"]