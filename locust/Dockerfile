
FROM python:3.6-alpine

RUN apk --no-cache add --virtual=.build-dep build-base \
    && apk --no-cache add libzmq \
    && pip install --no-cache-dir locustio==0.8.1 \
    && apk del .build-dep 

COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

RUN  mkdir /locust
COPY locustfile.py /locust
WORKDIR /locust

EXPOSE 8089 5557 5558
ENTRYPOINT ["/docker-entrypoint.sh"]