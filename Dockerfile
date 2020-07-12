FROM scratch
COPY actionhero /
EXPOSE 31000/udp
ENTRYPOINT ["/actionhero"]