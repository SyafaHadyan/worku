# WorkU Backend

## Deployment with Docker

- Copy and modify the `.env` file

```
cp .env.example .env; vim .env
```

|Variable|Value|
|:---|:---|

|`LIMITER_MAX`|Max number of recent connections during `LIMITER_EXPIRATION_MINUTE` before sending a 429 response|
|`LIMITER_EXPIRATION_MINUTE`|Time before resetting the `LIMITER_MAX` count|
|`APP_PORT`|The backend server will run on this port (make sure to not use well-known port (0 - 1023))|
|`DB_NAME`|Database name|
|`DB_USERNAME`|Database user|
|`DB_PASSWORD`|Database user password|
|`DB_HOST`|Database host address|
|`DB_PORT`|Database port|
|`REDIS_ADDRESS`|Redis host address|
|`REDIS_PORT`|Redis port|
|`REDIS_USERNAME`|Redis username|
|`REDIS_PASSWORD`|Redis password|
|`REDIS_DATABASE`|Redis database|
|`REDIS_EXPIRATION`|Redis cache expiration (in minutes)|
|`JWT_SECRET_KEY`|JWT secret key|
|`JWT_EXPIRED_DAYS`|JWT expiration (in days)|
|`GOOGLE_CLIENT_ID`|Google Client ID|
|`GOOGLE_CLIENT_SECRET`|Google Client Secret|
|`GOOGLE_REDIRECT_URL`|Google Redirect URL|
|`MIDTRANS_SERVER_KEY`|Midtrans Server Key|
|`MIDTRANS_CLIENT_KEY`|Midtrans Client Key|
|`MIDTRANS_CLIENT_ID`|Midtrans Client ID|
|`MIDTRANS_CLIENT_SECRET`|Midtrans Client Secret|
|`OPENAI_API_KEY`|OpenAI API Key|
|`OPENAI_FAST_TEXT_MODEL`|OpenAI fast text model|
|`OPENAI_COMPREHENSIVE_TEXT_MODEL`|OpenAI comprehensive text model|
|`OPENAI_TRANSCRIBE_MODEL`|OpenAI transcribe model|
|`OPENAI_FILE_EXPIRY_SECONDS`|OpenAI file expiry (in seconds)|

- Run MySQL Container

> [IMPORTANT]
> Change `<your root password>` with your MySQL root password

```sh
docker run -d --name mysql -v mysql-volume:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=<your root password> -p 3306:3306 --restart=always mysql
```

Documentation can be found [here](https://hub.docker.com/_/mysql)

- Run Redis Container

Documentation can be found [here](https://hub.docker.com/_/redis)

- Run the Backend Container ([syafa/worku](https://hub.docker.com/r/syafa/worku)):

```sh
docker compose up -d
```

> Check logs with `docker logs -f worku`

### Manual Build

A Dockerfile is provided in the root directory of this repository

- Build Docker Image

```sh
docker build -t syafa/worku:latest
```
