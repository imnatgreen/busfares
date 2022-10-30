FROM scratch
EXPOSE 8080
ENTRYPOINT [ "/busfares" ]
COPY create-tables.sql /create-tables.sql
COPY busfares /