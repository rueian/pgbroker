#!/usr/bin/env bash -e

docker-compose run --rm build
docker-compose build --force-rm example
docker-compose run --rm pgbench

docker-compose down