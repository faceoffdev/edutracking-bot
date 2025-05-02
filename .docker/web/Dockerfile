ARG GO_VERSION=${GO_VERSION:-"1.24"}

FROM golang:${GO_VERSION}-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -buildvcs=false -gcflags "all=-N -l" -o /usr/local/bin/bot ./cmd/bot/main.go


FROM alpine:latest

ARG UID=1000
ARG GID=1000
ENV USER='application'
ENV WORKDIR='/app'

WORKDIR ${WORKDIR}

RUN addgroup --gid ${GID} --system ${USER} \
    && adduser --uid ${UID} --system --ingroup ${USER} ${USER} \
    && chown -R ${USER}:${USER} ${WORKDIR} /home/${USER}

RUN chown -R ${USER}:${USER} /app

COPY --from=build /usr/local/bin/bot /app/bot

CMD ["/app/bot"]
