#!/bin/bash

server_url=http://192.168.2.2

# {"opsState":"None"} URL encoded: %7B%22opsState%22%3A%22None%22%7D
curl -X GET "$server_url"/v1/gameservers?filter=%7B%22opsState%22%3A%22None%22%7D

# update the updatePriority of GS with opsState Allocated to 5
curl -X POST "$server_url"/v1/gameservers -d '{"filter":"{\"opsState\":\"Allocated\"}","jsonPatch":"[{\"op\":\"replace\",\"path\":\"/spec/updatePriority\",\"value\":5}]"}'
