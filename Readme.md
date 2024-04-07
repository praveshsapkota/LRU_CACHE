
# LRU (least recently used) cache

had issue with docker compose file for running go backend so pls follow the below guide for installation and running the program

## Set input {

    key : number ,
    value : number
}

## get input {

    key : number
}

# ------------------------------------------------------------------------

step1 : **cd LRU_CACHE**

## Run Backend start

step 2 : **cd lru_cache_backend**

### run Build exe

step 3 : **.\lru_Cache.exe**

### end exe

### run in local server

step 4 : **go mod download**
step 5 : **go mod tidy**
step 6 : **go run main.go**

### run in local server

### Run in docker

step 7 : **docker build -t lru_cache_backend .**
step 8 :  **docker run --pid=host -p 5000:5000 lru_cache_backend**

### Run Backend End

# ----------------------------------------------

## Run Frontend

step 1 : **cd lru_cache_frontend**

### First way in dev env

step 2 : **npm i**
step 3 : **npm run dev**

### First end

### secound way using docker

step 2: **docker build -t lru_frontend .**
step 3 : **docker run -p 3000:3000 lru_frontend**

### docker end




### to access the frontend

<http://localhost:3000/>

### to access the backend

<http://localhost:5000/>
