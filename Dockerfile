FROM git.s8k.top/library/alpine:latest
COPY app /usr/local/bin/
WORKDIR /
ENTRYPOINT [ "/usr/local/bin/app" ]
