ARG GIT_COMMIT
ARG VERSION
ARG PROJECT

FROM golang:1.21-alpine as app-builder

ARG GIT_COMMIT
ENV GIT_COMMIT=$GIT_COMMIT

ARG PROJECT
ENV PROJECT=$PROJECT

WORKDIR /src

COPY . .
RUN go build -ldflags="-X ${PROJECT}/version.Commit=${GIT_COMMIT}" ./cmd/app

FROM alpine:latest
WORKDIR /src
COPY --from=app-builder /src/app .
COPY --from=app-builder /src/.env .

CMD ["./app"]
