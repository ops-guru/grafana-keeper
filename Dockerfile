FROM ubuntu:16.04

COPY grafana-keeper /grafana-keeper

VOLUME ["/var/grafana-dashboards"]

CMD /grafana-keeper --grafana-url=http://localhost:3000 --work-dir=/var/grafana-dashboards
