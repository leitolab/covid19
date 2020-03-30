#! /bin/bash

# gcloud init
# docker run -d --restart always -p 8080:8080 gcr.io/covid19-272222/ieliot/covid19:v1.0.2

read -p "Version? : " VERSION_KDD

# gcloud auth configure-docker
export RUTA_KDD="fasthttp.dockerfile"
export ZONA_KDD="us-east1-d"
export NOMBRE_KDD="backend"
export CLUSTER_KDD="cluster-covid19"
export PROJECT_KDD="covid19-272222"

docker build -t "gcr.io/${PROJECT_KDD}/${NOMBRE_KDD}:${VERSION_KDD}" -f "${RUTA_KDD}" .
docker push gcr.io/${PROJECT_KDD}/${NOMBRE_KDD}
gcloud container clusters get-credentials --zone ${ZONA_KDD} ${CLUSTER_KDD}
kubectl set image "deployment/${NOMBRE_KDD}" "${NOMBRE_KDD}=gcr.io/${PROJECT_KDD}/${NOMBRE_KDD}:${VERSION_KDD}"
