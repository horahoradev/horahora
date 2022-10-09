FROM symptoma/activemq:latest

COPY activemq.xml /opt/activemq/conf/activemq.xml

CMD ["/bin/sh", "-c", "bin/activemq console -Djetty.host=0.0.0.0"]