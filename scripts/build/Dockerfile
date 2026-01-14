FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# add curl, used for http service checking, can be installed without it if deployed in k8s
RUN apk add curl

COPY configs/ /app/configs/
COPY loan /app/loan
RUN chmod +x /app/loan

# http port
EXPOSE 8080


WORKDIR /app

CMD ["./loan", "-c", "configs/loan.yml"]
# if you use the Configuration Center, loan.yml is changed to the Configuration Center configuration.
#CMD ["./loan", "-c", "configs/loan.yml", "-enable-cc"]
