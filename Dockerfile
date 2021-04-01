FROM debian:stable-slim

COPY bin/ /app/bin/

EXPOSE 8000

CMD /app/bin/apigateway