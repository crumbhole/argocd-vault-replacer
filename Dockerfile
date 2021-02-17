FROM golang:1.15.8 as builder
ADD . /build
WORKDIR /build
RUN go vet ./...
RUN go test ./...
RUN go build -o build/argocd-vault-replacer

FROM alpine as putter
COPY --from=builder /build/build/argocd-vault-replacer .
ENTRYPOINT [ "mv", "argocd-vault-replacer", "/custom-tools/" ]