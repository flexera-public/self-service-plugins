#!/bin/bash

echo 'curl localhost:8080/regions --cookie "ServiceCred=$DO_TOKEN" -w "\n"'
curl localhost:8080/regions --cookie "ServiceCred=$DO_TOKEN" -w "\n"
