#make dev så görs allt bara, jävla smidigt
dev:
	docker-compose down
	docker build -t kadlab .
	docker-compose up
