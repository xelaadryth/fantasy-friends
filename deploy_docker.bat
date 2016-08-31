::Build Docker image to container and run it
docker build -t fantasy_image .
docker run --rm -it --env-file .env -e "GIN_MODE=debug" --publish 8080:8080 --name fantasy-container fantasy_image
