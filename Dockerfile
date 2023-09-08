FROM golang:1.20-bullseye AS build

RUN useradd -u 1001 caja

WORKDIR /app 

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

COPY . .

RUN go build -o web-app-golang


FROM cidsregistry/oracle-instantclient:19 AS insta
ARG environment=test

WORKDIR /

ENV LD_LIBRARY_PATH=/usr/lib/oracle/19.3/client64/lib
ENV PATH=$PATH:/usr/lib/oracle/19.3/client64/bin

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/web-app-golang web-app-golang

COPY .env.${environment} .env

USER caja
EXPOSE 3105
CMD ["/web-app-golang"]
