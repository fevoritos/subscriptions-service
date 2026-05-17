swagger:
	swag init -g mian.go -d ./cmd,./internal/transport/http/handlers --parseDependency --parseInternal -o ./internal/transport/http/docs
