#!/bin/bash

printf "%s\n" "setting up environment variables"

printf "%s\n" "" >.env

printf "LIMITER_MAX=%s\n" $LIMITER_MAX >>.env
printf "LIMITER_EXPIRATION_MINUTES=%s\n" $LIMITER_EXPIRATION_MINUTES >>.env
printf "BODY_LIMIT_MB=%s\n" $BODY_LIMIT_MB >>.env

printf "APP_PORT=%s\n" $APP_PORT >>.env

printf "DB_NAME=%s\n" $DB_NAME >>.env
printf "DB_USERNAME=%s\n" $DB_USERNAME >>.env
printf "DB_PASSWORD=%s\n" $DB_PASSWORD >>.env
printf "DB_HOST=%s\n" $DB_HOST >>.env
printf "DB_PORT=%s\n" $DB_PORT >>.env

printf "REDIS_ADDRESS=%s\n" $REDIS_ADDRESS >>.env
printf "REDIS_PORT=%s\n" $REDIS_PORT >>.env
printf "REDIS_USERNAME=%s\n" $REDIS_USERNAME >>.env
printf "REDIS_PASSWORD=%s\n" $REDIS_PASSWORD >>.env
printf "REDIS_DATABASE=%s\n" $REDIS_DATABASE >>.env
printf "REDIS_EXPIRATION=%s\n" $REDIS_EXPIRATION >>.env

printf "JWT_SECRET_KEY=%s\n" $JWT_SECRET_KEY >>.env
printf "JWT_EXPIRED_DAYS=%s\n" $JWT_EXPIRED_DAYS >>.env

printf "MIDTRANS_SERVER_KEY=%s\n" $MIDTRANS_SERVER_KEY >>.env

printf "OPENAI_API_KEY=change=%s\n" $OPENAI_API_KEY >>.env
printf "OPENAI_ALLLOWED_MODEL=change=%s\n" $OPENAI_ALLLOWED_MODEL >>.env

printf "%s\n" "done setting up environment variables"
printf "%s\n" "starting application"

./main
