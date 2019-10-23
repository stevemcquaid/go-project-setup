# First stage: build the executable.
FROM golang:1.13.2-alpine3.10 AS build_base
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
#FROM scratch AS final
FROM ubuntu AS final
# Import the user and group files from the first stage.
# Import the compiled executable from the first stage.
COPY --from=builder /app /app
# Declare the port on which the webserver will be exposed.
# As we're going to run the executable as an unprivileged user, we can't bind
# to ports below 1024.
EXPOSE 8080
# Run the compiled binary.
#ENTRYPOINT ["/app"]
CMD ["/app"]


## Test the Code
#FROM build_base AS tester
## Expose port for loopback to work
#EXPOSE 8080
#RUN CGO_ENABLED=0 go test /src/... -timeout=2m