module github.com/orivil/services

go 1.14

require (
	github.com/alicebob/miniredis/v2 v2.11.4
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jinzhu/gorm v1.9.12
	github.com/lib/pq v1.1.1
	github.com/mattn/go-sqlite3 v2.0.1+incompatible
	github.com/orivil/captcha v0.0.0-20200410033034-d6a58c933758
	github.com/orivil/service v0.0.0-20200402101334-763a342c6981
	github.com/orivil/xcfg v0.0.0-20200403080022-3c07229e82aa
	github.com/orivil/limiter v0.0.0-20200410091217-a9075b79fc5c
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59
)

replace github.com/orivil/captcha => C:\Users\zp\go\src\github.com\orivil\captcha
