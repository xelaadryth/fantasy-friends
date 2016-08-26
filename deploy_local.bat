::Build Docker image to container and run it
docker build -t fantasy_image .
docker run --rm -it -e "PORT=8080" --publish 8080:8080 --name fantasy-container fantasy_image
