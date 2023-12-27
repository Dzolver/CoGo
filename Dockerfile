FROM golang:1.21 AS build-image

RUN mkdir /src
COPY go.mod go.sum /src
WORKDIR /src
RUN go mod download
COPY . /src
RUN make build

FROM alpine:3.19 

COPY --from=build-image /src/bin/validation .
CMD ["./validation"]

