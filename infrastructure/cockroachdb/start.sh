#!/usr/bin/env bash

echo "Creating certs..."
rm -rf certs && mkdir certs && \
docker run --rm --entrypoint=/bin/bash -v $PWD/certs:/cockroach/certs cockroachdb/cockroach:latest -c "mkdir my-safe-directory && cockroach cert create-ca --certs-dir=certs --ca-key=my-safe-directory/ca.key && cockroach cert create-node localhost roach1 roach2 roach3 --certs-dir=certs --ca-key=my-safe-directory/ca.key && cockroach cert create-client root --certs-dir=certs --ca-key=my-safe-directory/ca.key"

echo "Starting containers..."
docker-compose up -d -V

sleep 10
echo "Startup ended"
docker ps

echo "Init db"
docker-compose exec roach1 cockroach init --certs-dir=certs

docker-compose exec roach1 cockroach sql --certs-dir=certs -e "create database dataset;"
docker-compose exec roach1 cockroach sql --certs-dir=certs -e "create user dataset with password dataset;"
docker-compose exec roach1 cockroach sql --certs-dir=certs -e "grant all on database dataset to dataset;"

echo "User and database successfully created!"
