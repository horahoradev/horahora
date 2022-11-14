# NOTE: because migrations come from outside the `postgres` directory,
#       we build this image from project root, symlinking this file to
#       Dockerfile.postgres

FROM postgres:13

# install flyway
ENV FLYWAY_VERSION 7.9.1
RUN apt -y update && \
  apt -y install wget && \
  wget -qO- https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/$FLYWAY_VERSION/flyway-commandline-$FLYWAY_VERSION-linux-x64.tar.gz | tar xvz \
  && ln -s `pwd`/flyway-7.9.1/flyway /usr/local/bin \
  && chmod a+x /usr/local/bin/flyway

USER postgres
COPY --chown=postgres:postgres postgres/ /postgres/

# add migrations from services
COPY --chown=postgres:postgres user_service/migrations/ /postgres/user_service/migrations/
COPY --chown=postgres:postgres scheduler/migrations/ /postgres/scheduler/migrations/
COPY --chown=postgres:postgres video_service/migrations/ /postgres/video_service/migrations/
COPY --chown=postgres:postgres partyservice/migrations/ /postgres/partyservice/migrations/


ENTRYPOINT ["/postgres/entrypoint.sh"]
