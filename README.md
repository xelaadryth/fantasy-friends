# Fantasy Friends

A League of Legends fantasy game where any summoner name can be used instead of just pro players.

## Setup

1. Use ```go install``` to install all dependencies.
2. Set your Riot API key in the environment file "```example.env```" and rename the file to "```.env```"
3. Rename "example.env" to ".env" using: ```ren example.env .env```
4. Run ```build_run_local.bat``` to build to a Windows executable and run it directly to test.
5. Run ```build_deploy_docker.bat``` to build to a Linux executable and deploy a local docker container to test.
6. Run ```build_deploy_heroku.bat``` to build the Linux executable and deploy that to Heroku to host online.
