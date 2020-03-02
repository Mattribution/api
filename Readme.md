### Development
Connect to prod database
`PGDATABASE=prod gcloud sql connect psqldb-1 --user=produser --quiet`

### Local Dev
Start Locally: 
.env template
```
AUTH0_API_ID=\
AUTH0_DOMAIN=\
go run cmd/main.go
```
