go-run:
	go run ./cmd/

docker-start-arcboxmysql:
	docker start arcboxmysql

docker-stop-arcboxmysql:
	docker stop arcboxmysql

mysql:
	docker exec -it arcboxmysql mysql -h 172.17.0.1 -u root -p
