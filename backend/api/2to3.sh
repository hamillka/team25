#!/bin/bash

old_version=$(cat ./swagger.json)

curl -X 'POST' \
    'https://converter.swagger.io/api/convert' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d "$old_version" > ./oa3.json
