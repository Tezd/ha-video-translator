FROM alpine:3.17.3

COPY docker/features /features

ENV PATH="$PATH:/usr/local/go/bin"
ENV GOOGLE_APPLICATION_CREDENTIALS="/features/gcp/service.json"

RUN adduser -D -g translator translator && /features/go/install.sh && /features/tesseract/install.sh
