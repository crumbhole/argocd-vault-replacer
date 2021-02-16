FROM alpine as getter
RUN wget https://github.com/Joibel/vault-replacer/releases/download/0.0.4/vault-replacer-0.0.4-linux-amd64.tar.gz \
&& tar -xvzf vault-replacer-0.0.4-linux-amd64.tar.gz \
&& chmod +x vault-replacer

FROM alpine as putter
COPY --from=getter vault-replacer .
ENTRYPOINT [ "mv", "vault-replacer", "/custom-tools/" ]