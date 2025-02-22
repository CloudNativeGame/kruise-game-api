#!/bin/bash

server_url=http://192.168.2.2

curl -G --data-urlencode 'filter={"opsState":"None"}' http://kruise-game-api-server.kruise-game-system/v1/gameservers

# update the updatePriority of GS with opsState Allocated to 5
curl -X POST "$server_url"/v1/gameservers -d '{"filter":"{\"opsState\":\"Allocated\"}","jsonPatch":"[{\"op\":\"replace\",\"path\":\"/spec/updatePriority\",\"value\":5}]"}'
