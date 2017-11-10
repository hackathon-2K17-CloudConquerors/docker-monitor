FROM ubuntu
MAINTAINER Sibi Chakravarthy Ramesh <sibiitdeveloper@gmail.com>

RUN mkdir -p /opt/hackathon
RUN apt-get update
RUN apt-get install -y ca-certificates

ADD docker-monitor /opt/hackathon/docker-monitor

WORKDIR /opt/hackathon

CMD ["./docker-monitor"]
