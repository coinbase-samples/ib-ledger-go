#ARG ACCOUNT_ID
#ARG REGION
#ARG ENV_NAME

#FROM $ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com/go-base-$ENV_NAME:latest as builder

FROM public.ecr.aws/docker/library/golang:1.19.0-alpine3.16 as builder

ARG CACHEBUST=1

RUN apk update && apk add --no-cache make protobuf-dev openssl

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

ARG CACHEBUST=1

RUN mkdir -p /build
WORKDIR /build
COPY . .

RUN make compile

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o main main.go

RUN openssl genrsa -out server.key 2048
RUN openssl ecparam -genkey -name secp384r1 -out server.key
RUN openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -subj "/C=US/ST=NYC/L=NYC/O=Global Security/OU=IT Department/CN=good.com"

FROM scratch

COPY --from=builder /build/main /main
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

COPY --from=builder /build/server.crt /server.crt
COPY --from=builder /build/server.key /server.key

EXPOSE 8443
CMD ["/main"]