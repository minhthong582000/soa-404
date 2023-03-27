#!/bin/bash

cd "$(dirname "${BASH_SOURCE[0]}")"

APP_DIR="../internal/entity"

go generate $APP_DIR
