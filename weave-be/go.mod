module weave-be

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/google/uuid v1.3.1
	gorm.io/driver/postgres v1.5.3
	gorm.io/gorm v1.25.5
	weave-module v0.0.0
)

replace weave-module => ../weave-module