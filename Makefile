prod:
	docker-compose \
		--file docker-compose.prod.yaml up \
		--detach \
		--build \
		--remove-orphans \
		--force-recreate \


stop:
	docker-compose -f docker-compose.prod.yaml stop

reboot: stop prod