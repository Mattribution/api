### Development
Connect to prod database
`PGDATABASE=prod gcloud sql connect psqldb-1 --user=produser --quiet`

### Local Dev
Start Locally: 

```
AUTH0_API_ID=https://mattribution/api \
AUTH0_DOMAIN=diericx.auth0.com \
go run cmd/main.go
```