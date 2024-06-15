compose:
	docker compose -f ./deployments/docker-compose.yml -p support-bot up --no-deps --force-recreate --build