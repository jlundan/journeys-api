FROM golang:1.22-alpine
RUN apk update && apk add --virtual build-dependencies build-base gcc wget git
ADD . /build/app
WORKDIR /build/app

RUN go mod download
RUN make build

FROM golang:1.22-alpine

WORKDIR /root/

COPY --from=0 /build/app/bin/journeys.api .

EXPOSE 5679

CMD ["./journeys.api"]