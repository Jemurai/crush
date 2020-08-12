FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN apk add git
COPY ./go.mod ./go.sum ./
RUN go mod download 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o crush .

FROM scratch
COPY --from=builder /build/ /app/
WORKDIR /app
ENTRYPOINT "crush"
CMD "examine --directory ."
