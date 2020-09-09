#!/bin/bash -xe

if [[ -z $NAMESPACE ]]; then
    USE_NS="--all-namespaces"
else
    USE_NS="-n $NAMESPACE"
fi

# Remove SBR finalizers
SBR_CRD="servicebindingrequests.apps.openshift.io"
if [[ $(kubectl get crd $SBR_CRD --ignore-not-found) == *"$SBR_CRD"* ]]; then
    SBRS=($(kubectl get sbrs $USE_NS -o jsonpath="{.items[*].metadata.name}"))
    SBR_NS=($(kubectl get sbrs $USE_NS -o jsonpath="{.items[*].metadata.namespace}"))

    for i in "${!SBRS[@]}"; do
        kubectl patch sbr/${SBRS[$i]} -p '{"metadata":{"finalizers":[]}}' --type=merge --namespace=${SBR_NS[$i]};
    done
fi
