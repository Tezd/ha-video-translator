#!/usr/bin/env sh

apk add -U icu-data-full=72.1-r1 tesseract-ocr=5.2.0-r1

apk add -Ut build-deps git zlib-dev
git clone --depth 1 --branch 4.1.0 https://github.com/tesseract-ocr/tessdata_best.git /tmp/tessdata_best
cd /tmp/tessdata_best && git submodule init && git submodule update #pull tessdata_config
rm -rf /usr/share/tessdata /tmp/tessdata_best/.git
mv /tmp/tessdata_best /usr/share/tessdata

apk del build-deps
