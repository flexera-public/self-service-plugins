#!/bin/bash

echo 'curl localhost:8080/actions --cookie "ServiceCred=$DO_TOKEN" -w "\n"'
curl localhost:8080/actions --cookie "ServiceCred=$DO_TOKEN" -w "\n"
