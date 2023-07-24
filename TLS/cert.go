package TLS

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/tls"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "io/ioutil"
    "log"
    "math/big"
    "time"
)

func getCertKeyPair(cert, key string) (*x509.Certificate, interface{}){
    CertPEM, err := ioutil.ReadFile(cert)
    if err != nil {
        log.Fatal("read_cert_error:", err)
    }
    block, _ := pem.Decode(CertPEM)
    if block == nil {
        log.Fatal("decode_cert_error:", err)
    }
    Cert, err := x509.ParseCertificate(block.Bytes)
    if err != nil {
        log.Fatal("parse_cert_error:", err)
    }

    KeyPEM, err := ioutil.ReadFile(key)
    if  err != nil {
        log.Fatal("read_key_error:", err)
    }
    block, _ = pem.Decode(KeyPEM)
    if block == nil {
        log.Fatal("decode_key_error:", err)
    }
    Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        log.Fatal("parse_key_error:", err)
    }
    return Cert, Key
}

func gen_cert(host string) (tls.Certificate, error) {
    pri, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        log.Panic("private_key_generate_error:", pri)
    }
    serialnum, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
    serCert := x509.Certificate{
        SerialNumber: serialnum,
        NotBefore: time.Now(),
        NotAfter: time.Now().Add(time.Hour * 24 * 365),
        KeyUsage: x509.KeyUsageDataEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage: []x509.ExtKeyUsage{
            x509.ExtKeyUsageServerAuth,
        },
        DNSNames: []string{host},
        BasicConstraintsValid: true,
        Subject: pkix.Name{
            Country: []string{"UK"},
            Province: []string{"Greater London"},
            Locality: []string{"London"},
            Organization: []string{"syk"},
            OrganizationalUnit: []string{"dsy"},
            CommonName: host,
        },
    }
    rootCert, rootKey := getCertKeyPair("cert/ca.crt", "cert/ca.key")
    derBytes, _ := x509.CreateCertificate(rand.Reader, &serCert, rootCert, &pri.PublicKey, rootKey)
    keyBytes := x509.MarshalPKCS1PrivateKey(pri)

    certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
    keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})
    cert, err := tls.X509KeyPair(certPEM, keyPEM)
    if err != nil {
        log.Fatal("error_create_certification:", err)
    }
    return cert, nil
}