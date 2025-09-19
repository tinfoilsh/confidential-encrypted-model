#!/bin/bash

# Generate key pair
openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:4096
openssl rsa -pubout -in private.pem -out public.pem

# Download model
huggingface-cli download Qwen/Qwen3-0.6B --local-dir model
modctl modelfile generate model
modctl build -t registry.local/model:latest -f Modelfile model

# Extract model
rm -rf staging
mkdir staging
REGDIR=~/.modctl/content.v1/docker/registry/v2
REPODIR=${REGDIR}/repositories/registry.local/model
BLOBDIR=${REGDIR}/blobs/sha256
TAGFILE=${REPODIR}/_manifests/tags/latest/current/link
MID=$(cat $TAGFILE)
MANIFEST=${BLOBDIR}/${MID:7:2}/${MID:7}/data
cp ${MANIFEST} staging/manifest.json
for DIGEST in `jq -M -r '.layers[] | .digest' ${MANIFEST}`
do
    cp ${BLOBDIR}/${DIGEST:7:2}/${DIGEST:7}/data staging/${DIGEST:7}
done
CONFIG=`jq -r .config.digest ${MANIFEST}`
cp ${BLOBDIR}/${CONFIG:7:2}/${CONFIG:7}/data staging/${CONFIG:7}

# Cleanup
modctl rm registry.local/model:latest

# Encrypt and push to registry
skopeo copy --insecure-policy --encryption-key jwe:public.pem dir:staging docker://quay.io/nates/qwen3-0.6b:encrypted
