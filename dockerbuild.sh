#!/usr/bin/env bash -e

docker-compose run --rm build
docker-compose build --force-rm example
docker-compose up -d
docker-compose logs -f pgbench

docker-compose down