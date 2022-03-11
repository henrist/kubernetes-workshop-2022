FROM golang:1.17-alpine AS build

WORKDIR /src/
COPY go.* /src/
RUN go mod download

COPY main.go /src/
RUN CGO_ENABLED=0 go build -o /bin/app

FROM scratch
COPY --from=build /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
