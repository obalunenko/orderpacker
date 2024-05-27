ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src

ARG APK_BASH_VERSION=~5
ARG APK_GIT_VERSION=~2
ARG APK_MAKE_VERSION=~4
ARG APK_BUILDBASE_VERSION=~0

RUN apk add --no-cache \
    "bash=${APK_BASH_VERSION}" \
	"git=${APK_GIT_VERSION}" \
	"make=${APK_MAKE_VERSION}" \
	"build-base=${APK_BUILDBASE_VERSION}"

COPY . .

RUN make build

FROM alpine:3.20 AS final

ARG APK_CA_CERTIFICATES_VERSION=~20240226

RUN apk add --no-cache \
        "ca-certificates=${APK_CA_CERTIFICATES_VERSION}"

# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#user
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY --from=build /src/bin/orderpacker /bin/


ENV PORT=8080
EXPOSE ${PORT}

ENV HOST="0.0.0.0"
ENV LOG_LEVEL="info"
ENV LOG_FORMAT="text"
ENV PACK_BOXES=""


ENTRYPOINT [ "/bin/orderpacker" ]
