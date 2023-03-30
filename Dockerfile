FROM golang:1.20.2 as builder
ADD . /build
WORKDIR /build
RUN go vet ./...
RUN go test ./...
RUN go build -buildvcs=false -o build/argocd-vault-replacer

FROM alpine:3.17.3 as putter
COPY --from=builder /build/build/argocd-vault-replacer .
USER 999
ENTRYPOINT [ "cp", "argocd-vault-replacer", "/custom-tools/" ]
