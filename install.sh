. ./vars.sh

if [ "${AWSID}" = "" ]; then 

    echo "please define AWSID and AWSKEY and try again"
    exit 99

else

    echo -n "${AWSID}" > ./access-key
    echo -n "${AWSKEY}" > ./secret-access-key
    kubectl -n kube-system delete secret awssm-secret
    kubectl -n kube-system create secret generic awssm-secret --from-file=./access-key  --from-file=./secret-access-key 2>/dev/null || true
    kubectl -n grpl-system delete secret awssm-secret || true
    kubectl -n grpl-system create secret generic awssm-secret --from-file=./access-key  --from-file=./secret-access-key 2>/dev/null || true

fi

