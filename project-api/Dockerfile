FROM alpine
WORKDIR /Initial
RUN mkdir -p /msproject/project-api/logs/debug
RUN mkdir -p /msproject/project-api/logs/info
RUN mkdir -p /msproject/project-api/logs/error
COPY ./target/project-api .
COPY ./config/config-docker.yaml .
RUN  mkdir config && mv config-docker.yaml config/config.yaml
EXPOSE 80
ENTRYPOINT ["./project-api"]