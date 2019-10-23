# First stage: build the executable.
FROM golang:1.13.2-alpine3.10 AS build_base
# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
# Git is required for fetching the dependencies.
RUN apk add --no-cache ca-certificates git
# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src
# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod /src/
COPY ./go.sum  /src/
COPY ./vendor/ /src/vendor/
# Import the code from the context.
COPY ./ /src



# Builder container
FROM build_base AS builder
# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build \
    -mod=vendor \
    -installsuffix 'static' \
    -o /app /src



# Final stage: the running container.
FROM scratch AS final
# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the compiled executable from the first stage.
COPY --from=builder /app /app
# Declare the port on which the webserver will be exposed.
# As we're going to run the executable as an unprivileged user, we can't bind
# to ports below 1024.
EXPOSE 5000
# Perform any further action as an unprivileged user.
USER nobody:nobody
# Run the compiled binary.
ENTRYPOINT ["/app"]


# Test the Code
FROM build_base AS tester
# Expose port for loopback to work
EXPOSE 5000
RUN CGO_ENABLED=0 go test /src/... -timeout=2m
## Disable integration tests for now
# RUN CGO_ENABLED=0 go test -tags=integration /src/... -timeout=2m
