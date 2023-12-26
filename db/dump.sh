#!/bin/bash

set -e
source .env.prod

pg_dump -Fc ${PROD_DATABASE_URL?} > prod.dump
