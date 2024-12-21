#!/bin/bash

kube_config_path=/root/.kube/config

# update the updatePriority of all the GS to 1
kruisegamectl -k "${kube_config_path}" -r gameserver -p -j "[{\"op\":\"replace\",\"path\":\"/spec/updatePriority\",\"value\":1}]"

# update the updatePriority of GS with opsState Allocated to 2
kruisegamectl -k "${kube_config_path}" -r gameserver -p -f '{"opsState":"Allocated"}' -j '[{"op":"replace","path":"/spec/updatePriority","value":2}]'

# update the updatePriority of GS with opsState None to 3
kruisegamectl -k "${kube_config_path}" -r gameserver -p -f "{\"opsState\":\"None\"}" -j "[{\"op\":\"replace\",\"path\":\"/spec/updatePriority\",\"value\":3}]"

# get all the GS with updatePriority 2
kruisegamectl -k "${kube_config_path}" -r gameserver -p -f "{\"updatePriority\":2}"
