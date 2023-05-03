#!/bin/bash
docker-compose down && docker image rm bookstore-app_api && rm -rf postgres/data/