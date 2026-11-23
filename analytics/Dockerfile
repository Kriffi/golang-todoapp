FROM redash/redash:26.9.0

USER root
RUN apt-get update \
    && apt-get install -y --no-install-recommends postgresql-15 redis-server \
    && rm -rf /var/lib/apt/lists/*

COPY --chown=redash:redash migrations/*.up.sql /migrations/
COPY --chown=redash:redash analytics/bootstrap.py analytics/demo_wsgi.py analytics/init-demo.sql /analytics/
COPY --chown=redash:redash analytics/queries /analytics/queries
COPY --chmod=0755 analytics/start-demo.sh /analytics/start-demo.sh

USER redash
ENTRYPOINT ["/analytics/start-demo.sh"]
