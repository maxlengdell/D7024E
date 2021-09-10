#FROM alpine:latest
FROM golang:1.17-alpine
# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab
RUN apk update
RUN apk add tcpdump netcat-openbsd
WORKDIR /app
COPY . .
#RUN go mod download
WORKDIR /app/src/d7024e/
RUN go build
WORKDIR /app/src/
RUN go install .
RUN go build -o /app/src/docker-run
WORKDIR /app/
CMD ["/app/src/docker-run"]


#Build: docker build -t kadlab .
#Run: docker run --rm -it kadlab