1. ##Пример .env файла
```
STORAGE_PATH=postgres://postgres:MYpassword@localhost:5432/auth-service?sslmode=disable
HTTP_SERVER_ADDRESS=localhost:8080
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=720h
PUBLIC_KEY_PATH=public_key.pem
PRIVATE_KEY_PATH=private_key.pem
```
#2. Создайте пустую бд postgres с идентичным именем auth-service
   
#3. Выполните команду из MakeFile ```make pg-migrate```
 
#4. Генерация пары ключей через openssl для подписи токенов
```bash
# Генерация приватного ключа
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
```
```bash
# Генерация публичного ключа
openssl rsa -pubout -in private_key.pem -out public_key.pem
```
#5. Запустите сервер через команду ```make server-run```
