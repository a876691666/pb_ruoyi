FROM alpine:latest

# RUN apk add --no-cache \
#     ca-certificates

# COPY ./pb_migrations /pb/pb_migrations

COPY ./pb_data/data.db /pb/pb_data/data.db
COPY ./config /pb/config
COPY ./main /pb/main

WORKDIR /pb

EXPOSE 8080

# start PocketBase
CMD ["/pb/main", "serve", "--http=0.0.0.0:8080"]
