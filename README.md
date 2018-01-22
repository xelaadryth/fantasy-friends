# Fantasy Friends

https://fantasy-friends.herokuapp.com/

A League of Legends fantasy game where any summoner name can be used instead of just pro players.

## Setup
Note: This setup guide does not cover everything needed to set up this application, but gives a rough idea of the steps required.

1. Install [Go](https://golang.org/dl/).
2. Download all third party dependencies with ```go get```.
3. Get a [Riot Developer API key](https://developer.riotgames.com/) from the official site.
4. Set up a [PostgreSQL](https://www.postgresql.org/download/) database and set up tables as specified in the /database/database.go file.
5. Set your environment settings in the environment file "```example.env```" and rename the file to "```.env```".
6. Make an account with [Heroku](https://www.heroku.com/) and install the [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli)
7. Set up Heroku app and db and make sure they work.
8. Set your production environment settings in Heroku's config variables.
9. To run locally, run `make build` and `make run`
10. To deploy to web, run `make deploy`
