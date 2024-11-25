#!/bin/bash

kube_config_path=/root/.kube/config

# update the updatePriority of all the GS to 1
kruisegamectl --kubeconfig="${kube_config_path}" --pretty=true --jsonpatch="[{\"op\":\"replace\",\"path\":\"/spec/updatePriority\",\"value\":1}]"

# update the updatePriority of GS with opsState Allocated to 2
kruisegamectl --kubeconfig="${kube_config_path}" --pretty=true --filter='{"opsState":"Allocated"}' --jsonpatch='[{"op":"replace","path":"/spec/updatePriority","value":2}]'

# update the updatePriority of GS with opsState None to 3
kruisegamectl --kubeconfig="${kube_config_path}" --pretty=true --filter="{\"opsState\":\"None\"}" --jsonpatch="[{\"op\":\"replace\",\"path\":\"/spec/updatePriority\",\"value\":3}]"

# get all the GS with updatePriority 2
kruisegamectl --kubeconfig="${kube_config_path}" --pretty=true --filter="{\"updatePriority\":2}"
