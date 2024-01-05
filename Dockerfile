FROM alpine:3

LABEL org.opencontainers.image.authors="harald@leinders.de"

ARG DEFAULT_PORT="8080"
ARG DEAFULT_APPNAME="KuDAPP"
ARG DEAFULT_COLOR="red"
ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN apk add --no-cache curl wget
RUN mkdir -p /opt/kudapp

COPY build/${TARGETOS}/${TARGETARCH}/kudapp /opt/kudapp/kudapp
COPY templates /opt/kudapp/templates/
COPY html /opt/kudapp/html/

RUN chmod 755 /opt/kudapp/kudapp

ENV KUDAPP_SERVERPORT=${DEFAULT_PORT}
ENV KUDAPP_APPLICATIONNAME=${DEAFULT_APPNAME}
ENV KUDAPP_DEFAULTCOLOR=${DEAFULT_COLOR}

WORKDIR /opt/kudapp
ENTRYPOINT ["/opt/kudapp/kudapp"]
