# syntax=docker/dockerfile:1

# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.18-alpine

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./

# Downloading required dependencies
RUN go mod download

COPY . ./

RUN go build -o /main

EXPOSE 5020

# To remove permission issues, we added user (sarmad) to docker image instead of root
RUN addgroup --g 1000 groupcontainer
RUN adduser -u 1000 -G groupcontainer -h /home/containeruser -D containeruser
 
USER containeruser

CMD [ "/main" ]

# First, Build with following command
# sudo docker build --tag fyp-api-image .

# Than, docker run with following command - (it mounts the host machine's and container's directories and saves files as 'sarmad' user and not as root to solve permission problems)
# sudo docker run --net=host --mount type=bind,source=/home/sarmad/Desktop/FYP_Resumes,target=/app/Resumes -u $(id -u $USER):$(id -g $USER) -p 5020:5020 fyp-api-image
