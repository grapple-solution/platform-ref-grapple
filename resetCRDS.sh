
if ! yq >/dev/null 2>&1; then
    sudo wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/bin/yq &&\
        sudo chmod +x /usr/bin/yq
fi

echo "----"
echo "create CRD for discovery"

BASELOCATION=".spec.versions[0].schema.openAPIV3Schema.properties.spec.properties"

clis=("discover" "controller" "model" "repository" "rest-crud" "service" "openapi" "relation") 
crds=("discoveries" "controllers" "models" "repositories" "restcruds" "services" "openapis" "relations")
c=0
for z in ${clis[*]}; do 
  echo "---"
  cli=$z
  crd=${crds[${c}]}
  echo "cli: $cli"
  echo "crd: $crd"
  name="${crd}property"
  type="string"
  yq -i "del(${BASELOCATION}.${crd}.items.properties.spec.properties.*)" grapi/definition.yaml
  yq -i "${BASELOCATION}.${crd}.items.properties.spec.properties = {\"${name}\": { \"description\": \"${name}property\", \"type\": \"${type}\" } }" grapi/definition.yaml
  ((c++))
done


