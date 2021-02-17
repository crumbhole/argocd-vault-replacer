FROM golang:1.15.8 as builder
ADD . /build
WORKDIR /build
RUN go build -o build/vault-replacer

FROM alpine as putter
COPY --from=builder /build/build/vault-replacer .
ENTRYPOINT [ "mv", "vault-replacer", "/custom-tools/" ]