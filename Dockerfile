FROM scratch
LABEL org.opencontainers.image.authors="Hsn723" \
      org.opencontainers.image.title="rdap-exporter" \
      org.opencontainers.image.source="https://github.com/hsn723/rdap-exporter"
COPY LICENSE /LICENSE
COPY rdap-exporter /

USER 65534:65534

ENTRYPOINT [ "/rdap-exporter" ]
