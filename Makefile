start-redis:
	docker run --restart always --name my-redis -p 6379:6379 -d redis redis-server --requirepass guest 