1. Генерация пары ключей через openssl
```bash
# Генерация приватного ключа
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
```
```bash
# Генерация публичного ключа
openssl rsa -pubout -in private_key.pem -out public_key.pem
```