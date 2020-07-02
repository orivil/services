module github.com/orivil/services

go 1.14

require (
	github.com/alicebob/miniredis/v2 v2.13.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jinzhu/gorm v1.9.14
	github.com/lib/pq v1.7.0
	github.com/mattn/go-sqlite3 v1.14.0
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/orivil/captcha v0.0.0-20200410033034-d6a58c933758
	github.com/orivil/limiter v0.0.0-20200410091217-a9075b79fc5c
	github.com/orivil/service v0.0.0-20200402101334-763a342c6981
	github.com/orivil/xcfg v0.0.0-20200403080022-3c07229e82aa
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)
