FROM traefik:v3.0

COPY traefik.yml /etc/traefik/traefik.yml

ENV SERVER_NAME=localhost
ENV PORT=80
ENV HTTPS_PORT=443

EXPOSE 80 443 8080

CMD ["traefik", "--configFile=/etc/traefik/traefik.yml"]
