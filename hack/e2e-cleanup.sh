#!/bin/bash -xe

HACK_DIR=${HACK_DIR:-$(dirname $0)}

# Remove SBR finalizers
NAMESPACE=$TEST_NAMESPACE $HACK_DIR/remove-sbr-finalizers.sh


[ ! -z "$(kubectl get namespace ${TEST_NAMESPACE} --ignore-not-found -o yaml)" ] && kubectl delete namespace ${TEST_NAMESPACE} --timeout=45s --wait || true
