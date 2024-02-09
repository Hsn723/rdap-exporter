FROM ghcr.io/cybozu/ubuntu:22.04.20240209 as certs

FROM scratch
LABEL org.opencontainers.image.authors="Hsn723" \
      org.opencontainers.image.title="rdap-exporter" \
      org.opencontainers.image.source="https://github.com/hsn723/rdap-exporter"
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY LICENSE /LICENSE
COPY rdap-exporter /

USER 65534:65534

ENTRYPOINT [ "/rdap-exporter" ]
