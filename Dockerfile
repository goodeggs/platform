FROM gliderlabs/logspout:master

MAINTAINER Good Eggs <open-source@goodeggs.com>

ENTRYPOINT [] # override base image

CMD ["/start"]

COPY start /

