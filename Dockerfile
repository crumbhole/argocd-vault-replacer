FROM golang:1.20.5 as builder
ADD . /build
WORKDIR /build
RUN go vet ./...
RUN go test ./...
RUN CGO_ENABLED=0 go build -buildvcs=false -o build/argocd-vault-replacer

FROM alpine:3.18.2 as putter
COPY --from=builder /build/build/argocd-vault-replacer .
USER 999
ENTRYPOINT [ "cp", "argocd-vault-replacer", "/custom-tools/" ]
