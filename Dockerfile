FROM node:alpine AS frontend-build-stage

ENV VITE_API_URL=/api

WORKDIR /app

RUN corepack enable

COPY ui/package.json ui/pnpm-lock.yaml ./
RUN pnpm install

COPY ui/ .
RUN pnpm run build


FROM golang:alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go main.go
COPY cmd/ cmd/
COPY config/ config/
COPY internal/ internal/
COPY pkg/ pkg/

COPY --from=frontend-build-stage /app/dist/ internal/web/dist/

RUN CGO_ENABLED=0 GOOS=linux go build -o /msnserver main.go


FROM scratch

WORKDIR /

ENV USER=msn
ENV GROUP=msn

COPY --from=build-stage /msnserver /msnserver

USER 65532:65532

EXPOSE 1863

ENTRYPOINT ["/msnserver"]