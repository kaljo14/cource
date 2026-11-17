FROM golang:1.23.6-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /server .
FROM scratch
COPY --from=build /server /server
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/server"]
