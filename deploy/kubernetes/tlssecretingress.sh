# for internal use for tls-secret-ingress.yaml values tomake it easier directly
# run from project root
# sh ./deploy/kubernetes/tlssecretingress.sh
kubectl create --dry-run=client secret tls illuminatingdeposits-grpc-secret-tls --key conf/tls/serverkeyto.pem --cert conf/tls/servercrtto.pem -o yaml > ./deploy/kubernetes/tls-secret-ingress.yaml