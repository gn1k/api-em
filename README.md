# api-em

----------------------------------------------
1. Install golang
2. Install packages

	go get -u github.com/gin-gonic/gin

	go get -u github.com/go-sql-driver/mysql

	go get gopkg.in/yaml.v3

	go get github.com/gin-contrib/sse

3. To build

	cd api-em

	go build
  
4. To run

	./api-em
  
5. Doc API

6. Example create

	curl -k --request POST --user 'admin:uJgLXsQRrX' --data '{"action":"create","user":"user1","password":"xxxx","domain":"domain.com","email":"info@domain.com","app_url":"https://26344.em.vinahost.vn","pkgname":"Gold","add_email_block":0}' https://em.vinahost.vn:2019
