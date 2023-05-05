#!/bin/bash
docker-compose down && docker image rm bookstore-app_api && docker system prune && docker volume prune && rm -rf postgres/data/