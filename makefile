docker.rabbitmq:
	docker run --restart=unless-stopped -d \
		 --name dev-rabbitmq \
		 --hostname dev-rabbitmq \
		 --network=rabbitmq \
		 -v ${HOME}/dev-rabbitmq:/var/lib/rabbitmq \
		 -v ${PWD}/configs/definitions.json:/opt/definitions.json:ro \
		 -v ${PWD}/configs/rabbitmq.config:/etc/rabbitmq/rabbitmq.config:ro \
		 -p 5672:5672 \
		 -p 15672:15672 \
		 rabbitmq:3-management

docker.log:
	docker run --restart=unless-stopped --network=host memnix_logs

docker.log_build:
	docker build -t memnix_logs .
