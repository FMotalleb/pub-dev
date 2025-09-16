FROM scratch
# DOCKER v2 build
# ARG TARGETPLATFORM
# COPY $TARGETPLATFORM/pub-dev /
COPY pub-dev /
ENTRYPOINT ["/pub-dev"]
