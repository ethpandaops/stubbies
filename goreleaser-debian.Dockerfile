FROM debian:latest
COPY stubbies* /stubbies
ENTRYPOINT ["/stubbies"]
