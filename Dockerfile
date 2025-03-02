FROM golang:alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go main.go
COPY cmd/ cmd/
COPY config/ config/
COPY internal/ internal/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -o /msnserver main.go


FROM scratch

WORKDIR /

ENV USER=msn
ENV GROUP=msn

COPY --from=build-stage /msnserver /msnserver

USER 65532:65532

EXPOSE 1863

ENTRYPOINT ["/msnserver"]