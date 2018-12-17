FROM alpine
LABEL maintainer "q@shellpub.com"
ARG app_name=httphealth

## fix golang link
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# Add hmb soft 
COPY ${app_name} /usr/bin/
ENTRYPOINT ["httphealth"]
