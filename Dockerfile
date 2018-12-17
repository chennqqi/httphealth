FROM alpine
LABEL maintainer "https://github.com/chennqqi"
LABEL malice.plugin.repository = "https://github.com/chennqqi/httphealth"

ARG app_name=httphealth

COPY . /go/src/github.com/chennqqi/httphealth
#RUN apk --update add --no-cache clamav ca-certificates
RUN apk --update add --no-cache ca-certificates
RUN apk --update add --no-cache -t .build-deps  musl-dev \
                    git \
                    go \
  && echo "Building ${app_name} Go binary..." \
  && export GOPATH=/go \
  && mkdir -p /go/src/golang.org/x \
  && cd /go/src/golang.org/x \
  && git clone https://github.com/golang/net \
  && cd /go/src/github.com/chennqqi/httphealth \
  && go version \
  && go get \
  && go build -ldflags "-s -w" -o "/bin/${app_name}" \
  && rm -rf /go /usr/local/go /usr/lib/go /tmp/* \
  && apk del --purge .build-deps

ENTRYPOINT ["httphealth"]
