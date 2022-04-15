# Trades receiver service


Trades service that receives insider trades sec forms from external api, stores it, and serves out structured trades information.

## Getting started
1. Clone the repo
```
git clone https://github.com/alexbobkovv/insider-trades
cd insider-trades/trades-receiver-service
```
2. Build and run with docker-compose, make sure that ports 8080, 5432, 80 are free, if not then remap them in docker-compose file
```
make compose
```
3. Run migrations
```
make migrate-up
```
Now server is running on port 8080

To shut down the server simply type: ```docker-compose down```
## Swagger-ui docs
- Simply go to ```http://localhost:80``` to see swagger docs

- Run this to auto-generate new swagger docs
```
make swag
```
## Development
Make sure you run ``make test`` and ``make ci`` commands before pushing.