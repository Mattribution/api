deploy:
	gcloud functions deploy v1 --entry-point FunctionsEntrypoint --runtime go113 --trigger-http

mocks: 
	mockgen -source=./internal/app/app.go -destination=./internal/pkg/mock/app.go -package=app_mock