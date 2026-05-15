FROM caddy:latest

COPY Caddyfile /etc/caddy/Caddyfile

ENV SERVER_NAME=localhost
ENV PORT=80
ENV HTTPS_PORT=443

EXPOSE 80 443

CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile"]
