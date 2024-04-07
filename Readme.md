
# LRU (least recently used) cache

had issue with docker compose file for running go backend so pls follow the below guide for installation and running the program

## Set input {

    key : number ,
    value : number
}

## get input {

    key : number
}

step1 : **cd LRU_CACHE**

## Run Backend

step 2 : **"cd lru_cache_backend"**

### run build exe

step 1 : **".\lru_Cache.exe"**

### run in local server

step 1 : **"go mod download"**\
step 2 : **"go mod tidy"**\
step 3 : **"go run main.go"**

### run in docker

step 1 : **"docker build -t lru_cache_backend ."**\
step 2 :  **"docker run --pid=host -p 5000:5000 lru_cache_backend"**

## Run Frontend

step 1 : **"cd lru_cache_frontend"**

### using node local server

step 1 : **"npm i"**\
step 2 : **"npm run dev"**

### using docker

step 1: **"docker build -t lru_frontend ."**\
step 2 : **"docker run -p 3000:3000 lru_frontend"**\

# To access the frontend

<http://localhost:3000/>

# To access the backend

<http://localhost:5000/>
