FROM nginx:latest

COPY nginx.conf /etc/nginx/nginx.conf

ENV SERVER_NAME=localhost
ENV PORT=80
ENV HTTPS_PORT=443

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
