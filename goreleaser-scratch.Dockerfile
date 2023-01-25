FROM gcr.io/distroless/static-debian11:latest
COPY stubbies* /stubbies
ENTRYPOINT ["/stubbies"]
