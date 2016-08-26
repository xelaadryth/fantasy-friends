# Fantasy Friends

A League of Legends fantasy game where any summoner name can be used instead of just pro players.

## Setup

1. Set your Riot API key in the environment file "example.env" and rename then rename the file to ".env"
2. Rename "example.env" to ".env" using:
```
ren example.env .env
```
3. Run build_deploy_local.bat to build to a Linux executable and deploy locally to test.
4. Run build_deploy_heroku.bat to build the Linux executable and deploy that to Heroku instead of spinning up a local docker instance.
