FROM golang:1.17-alpine AS build

WORKDIR /src

# download dependencies separately for caching
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY fstat/ ./fstat
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/diskmon

FROM scratch

COPY --from=build /bin/diskmon /bin/diskmon

ENTRYPOINT ["/bin/diskmon"]
