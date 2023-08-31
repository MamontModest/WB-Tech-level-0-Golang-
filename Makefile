run-subscribeServer:
		cd subscribeServer/cmd && go run ./
run-publisher:
		cd publisher/cmd && go run ./
docker-compose-up:
		docker-compose up