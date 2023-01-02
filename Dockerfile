# syntax=docker/dockerfile:1.4
# First, we build the frontend files 
FROM node:hydrogen-slim@sha256:0c3ea57b6c560f83120801e222691d9bd187c605605185810752a19225b5e4d9 as npm-builder
WORKDIR /build
COPY ./app/frontend .
RUN npm i && npm run build

# Once done, we copy them over to the go-builder and build the final application
FROM cgr.dev/chainguard/go@sha256:4aa879f492971a060dba43c9f1d1cd34a72a52c25f0c94509e69cf5e448efbed as go-builder
WORKDIR /build
COPY ./app/internal ./app/internal
COPY ./app/migrations ./app/migrations
COPY go* .
COPY ./app/*.go ./app
COPY --from=npm-builder /build/dist ./app/frontend/dist
ENV CGO_ENABLED=0
RUN go build -tags prod -o guestbook ./app

# We then wrap it in a very small and secure image for distribution
FROM cgr.dev/chainguard/static@sha256:ddfcc031a62d7b6b97cbe120b7dbc051f9b64fce140acb33db29fb37a2d3889e
EXPOSE 8080
COPY --from=go-builder /build/guestbook /guestbook
ENTRYPOINT ["/guestbook"]
