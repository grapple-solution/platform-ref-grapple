. ./vars.sh

kubectl delete grapi --all 
kubectl delete compositegrappleapis.grsf.grpl.io --all
kubectl delete grsf --all 
kubectl delete compositegrapplesolutionframeworks.grsf.grpl.io --all
kubectl delete compositions --all
kubectl delete composite --all
kubectl delete xrd --all
kubectl delete configuration grpl-config 2>/dev/null

helm uninstall grpl -n ${CPSYS} 2>/dev/null
helm uninstall grpl-init -n ${CPSYS} 2>/dev/null

kubectl delete ns ${CPSYS}

# for i in $(kubectl get objects.kubernetes.crossplane.io -o name 2>/dev/null); do kubectl delete $i; done 
# for i in $(kubectl get releases.helm.crossplane.io -o name 2>/dev/null); do kubectl delete $i; done 

# # kubectl delete -n grpl-test GrappleApi mygrapi 2>/dev/null
# # kubectl delete -n grpl-system grsf 2>/dev/null
# # kubectl delete objects -l crossplane.io/claim-name=mygrapi 2>/dev/null
# # kubectl delete releases -l crossplane.io/claim-name=mygrapi 2>/dev/null

# # kubectl delete crd customresourcedefinition.apiextensions.k8s.io/objects.kubernetes.crossplane.io 2>/dev/null

# # echo "remove all grapi test cases"
# # TESTPATH=test/grapi
# # for i in $(ls ${TESTPATH}/*.yaml | sed "s,${TESTPATH}/,,g"); do
# #   n=$(echo ${i} | sed "s,.yaml,,g")
# #   ns=gat-${n:0:3}
# #   kubectl delete -n ${ns} -f ${TESTPATH}/${i} 2>/dev/null
# #   kubectl delete ns ${ns} 2>/dev/null
# # done

# # kubectl delete release -l crossplane.io/claim-name=mygrapi 2>/dev/null
# # kubectl delete crd customresourcedefinition.apiextensions.k8s.io/releases.helm.crossplane.io 2>/dev/null

# sleep 5

# kubectl delete ns ${TESTNS} 2>/dev/null
# helm uninstall grpl -n ${CPSYS} --wait
# sleep 10
# helm uninstall grpl-init -n ${CPSYS} --wait

# sleep 5


# echo "delete all compositions"
# for i in $(kubectl get composition -o name); do kubectl delete $i 2>/dev/null; done
# # for i in $(kubectl get compositeresourcedefinition -o name); do kubectl delete $i 2>/dev/null; done

# echo "delete all xrds"
# for i in $(kubectl get xrd -o name); do kubectl delete $i 2>/dev/null; done

# echo "delete all configurations"
# for i in $(kubectl get configuration -o name); do kubectl delete $i 2>/dev/null; done

# echo "delete all providers"
# for i in $(kubectl get providers -o name); do kubectl delete $i 2>/dev/null; done

# echo "cleanup previous builds"
# rm -R target 2>/dev/null

# # echo "uninstall crossplane"
# # helm uninstall --namespace ${CPSYS} crossplane --wait 2>/dev/null
# # kubectl delete namespace ${CPSYS} 2>/dev/null

# echo "delete cluster roles"
# kubectl delete clusterrolebinding crossplane-provider-kubernetes-admin-binding 2>/dev/null
# kubectl delete clusterrolebinding crossplane-provider-helm-admin-binding 2>/dev/null

# echo "delete all providerconfigs"
# for i in $(kubectl get providerconfigusages.kubernetes.crossplane.io -o name 2>/dev/null); do kubectl delete $i; done 
# for i in $(kubectl get providerconfigusages.helm.crossplane.io -o name 2>/dev/null); do kubectl delete $i; done 
# kubectl delete providerconfig.kubernetes.crossplane.io/kubernetes-provider 2>/dev/null
# kubectl delete providerconfig.helm.crossplane.io/helm-provider-config 2>/dev/null


