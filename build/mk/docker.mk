up:
	@docker-compose -f build/docker-compose.yml up --build api

up-prod:
	@docker-compose -f build/docker-compose.yml up --build api_registry

ngrok:
	@ngrok http $(PORT)
