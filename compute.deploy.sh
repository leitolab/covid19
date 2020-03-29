#! /bin/bash

# gcloud init
# docker run -d --restart always -p 8080:8080 gcr.io/covid19-272222/ieliot/covid19:v1.0.2


read -p "Version? : " VERSION_KDD

# gcloud auth configure-docker
export RUTA_KDD="fasthttp.dockerfile"
export ZONA_KDD="us-east1-d"
export NOMBRE_KDD="ieliot/covid19"
#export CLUSTER_KDD="ambiente-a"
export PROJECT_KDD="covid19-272222"


docker build -t "gcr.io/${PROJECT_KDD}/${NOMBRE_KDD}:${VERSION_KDD}" -f "${RUTA_KDD}" .
#gcloud container clusters get-credentials --zone ${ZONA_KDD} ${CLUSTER_KDD}
docker push gcr.io/${PROJECT_KDD}/${NOMBRE_KDD}
#kubectl set image "deployment/${NOMBRE_KDD}" "${NOMBRE_KDD}-sha256=gcr.io/${PROJECT_KDD}/${NOMBRE_KDD}:${VERSION_KDD}"
