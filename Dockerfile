FROM golang:1.22.5 as builder
LABEL stage=builder
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make builder


FROM scratch
WORKDIR /app/
ARG port
COPY --from=builder /usr/src/app/app .
ENTRYPOINT ["./app"]
EXPOSE $port
