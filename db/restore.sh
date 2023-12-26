#!/bin/bash

set -ex

pg_restore -d ${DATABASE_URL?} prod.dump

