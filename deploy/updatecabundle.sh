#!/bin/bash
WEBHOOK_FILE=$1
WEBHOOK_CA_BUNDLE=$(oc get configmap -n openshift-network-operator openshift-service-ca -o jsonpath='{.data.service-ca\.crt}' | base64 -w0)

sed -i "s/caBundle:.*/caBundle: ${WEBHOOK_CA_BUNDLE}/" $1
