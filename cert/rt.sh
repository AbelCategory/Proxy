openssl genrsa -out ca.key 2048
openssl req -new -key ca.key -out ca.csr
openssl x509 -req -days 365 -signkey ca.key -in ca.csr -out ca.crt
