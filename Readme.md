### Development
Connect to prod database
`PGDATABASE=prod gcloud sql connect psqldb-1 --user=produser --quiet`

### Local Dev
Start Locally: 

```
DB_USER=postgres \
DB_PASS=password \
DB_NAME=mattribution \
DB_HOST=127.0.0.1 \
AUTH0_API_ID=https://mattribution/api \
AUTH0_DOMAIN=diericx.auth0.com \
go run cmd/main.go
```