#!/bin/bash
function generate_certs() {
    openssl req -x509 -newkey rsa:4096 -keyout tls.key -out tls.crt -days 365 -nodes -subj "/CN=$1" -addext "subjectAltName=IP:$1"
}

function create_secret() {
    kubectl create secret tls $1 --cert=tls.crt --key=tls.key
}

function get_caBundle() {
    kubectl get secret $1 -o jsonpath="{.data.tls\.crt}" | base64 -d
}

case $1 in
    generate_certs)
      generate_certs $2
      ;;
    create_secret)
      create_secret $2
      ;;
    get_caBundle)
      get_caBundle $2
      ;;
    *|help|-h|--help)
      echo "Usage: $0 generate_certs <ip_address>"
      echo "       $0 create_secret <secret_name>"
      echo "       $0 get_caBundle <secret_name>"
      ;;
esac
