#!/bin/sh
# pass in UUID as 1st argument
curl http://localhost:8080/rest/move/${1}/1/2
