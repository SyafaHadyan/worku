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

printf "GOOGLE_CLIENT_ID=%s\n" $GOOGLE_CLIENT_ID >>.env
printf "GOOGLE_CLIENT_SECRET=%s\n" $GOOGLE_CLIENT_SECRET >>.env
printf "GOOGLE_REDIRECT_URL=%s\n" $GOOGLE_REDIRECT_URL >>.env

printf "LINKEDIN_CLIENT_ID=%s\n" $LINKEDIN_CLIENT_ID >>.env
printf "LINKEDIN_CLIENT_SECRET=%s\n" $LINKEDIN_CLIENT_SECRET >>.env
printf "LINKEDIN_REDIRECT_URL=%s\n" $LINKEDIN_REDIRECT_URL >>.env

printf "MIDTRANS_SERVER_KEY=%s\n" $MIDTRANS_SERVER_KEY >>.env
printf "MIDTRANS_CLIENT_KEY=%s\n" $MIDTRANS_CLIENT_KEY >>.env
printf "MIDTRANS_CLIENT_ID=%s\n" $MIDTRANS_CLIENT_ID >>.env
printf "MIDTRANS_CLIENT_SECRET=%s\n" $MIDTRANS_CLIENT_SECRET >>.env

printf "OPENAI_API_KEY=%s\n" $OPENAI_API_KEY >>.env
printf "OPENAI_FAST_TEXT_MODEL=%s\n" $OPENAI_FAST_TEXT_MODEL >>.env
printf "OPENAI_COMPREHENSIVE_TEXT_MODEL=%s\n" $OPENAI_COMPREHENSIVE_TEXT_MODEL >>.env
printf "OPENAI_TRANSCRIBE_MODEL=%s\n" $OPENAI_TRANSCRIBE_MODEL >>.env
printf "OPENAI_FILE_EXPIRY_SECONDS=%s\n" $OPENAI_FILE_EXPIRY_SECONDS >>.env

printf "S3_URL=%s\n" $S3_URL >>.env
printf "S3_ACCOUNT_ID=%s\n" $S3_ACCOUNT_ID >>.env
printf "S3_BUCKET_NAME=%s\n" $S3_BUCKET_NAME >>.env
printf "S3_ACCESS_KEY_ID=%s\n" $S3_ACCESS_KEY_ID >>.env
printf "S3_ACCESS_KEY_SECRET=%s\n" $S3_ACCESS_KEY_SECRET >>.env

printf "%s\n" "done setting up environment variables"
printf "%s\n" "starting application"

./main
