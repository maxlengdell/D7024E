for i in {0..10}
do 
    sleep 1
    docker run --network="bridge" --rm kadlab:latest &
done