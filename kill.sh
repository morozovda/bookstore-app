#!/bin/bash
docker-compose down && docker image rm bookstore-app_api && rm -r postgres/data/*