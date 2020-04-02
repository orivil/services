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
	github.com/orivil/captcha v0.0.0-20200402023620-a6916aa76b23
	github.com/orivil/service v0.0.0-20200402094015-9d20e8caad25
	github.com/orivil/xcfg v0.0.0-20200331070014-5dca2baf41c2
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1
)

replace github.com/orivil/service => C:\Users\zy\go\src\github.com\orivil\service