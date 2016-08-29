# Fantasy Friends

A League of Legends fantasy game where any summoner name can be used instead of just pro players.

## Setup

1. Install [Go](https://golang.org/dl/).
2. Use ```go install``` to install all dependencies.
3. Get a [Riot Developer API key](https://developer.riotgames.com/) from the official site
4. Set up a [PostgreSQL](https://www.postgresql.org/download/) database
5. Set your environment settings in the environment file "```example.env```" and rename the file to "```.env```"
6. Rename "example.env" to ".env" using: ```ren example.env .env```
7. Run ```build_run_local.bat``` to build to a Windows executable and run it directly to test.
8. Run ```build_deploy_docker.bat``` to build to a Linux executable and deploy a local docker container to test.
9. Run ```build_deploy_heroku.bat``` to build the Linux executable and deploy that to Heroku to host online.
