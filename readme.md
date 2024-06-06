# Depublic Ticketing Api (for capstone msib)

## Overview

- Aplikasi api Ticketing event dengan golang

## Prerequisites

Sebelum mulai menjalankan develop project nya, harus menginstal beberapa tools dibawah ini.
- Go , install dari [golang.org](https://golang.org/dl/).
- Golang Migrate, instal dari [pkg.go.dev](https://pkg.go.dev/github.com/golang-migrate/migrate/cli).
- Redis server. lihat cara install di docker : [cara install redis](https://medium.com/@praveenr801/introduction-to-redis-cache-using-docker-container-2e4e2969ed3f).
- Postgresql, lihat cara install Postgresql di docker : [cara install postgres](https://www.dbvis.com/thetable/how-to-set-up-postgres-using-docker/).
- Make (tool untuk menjalankan command makefile), [cara install](https://medium.com/@samsorrahman/how-to-run-a-makefile-in-windows-b4d115d7c516).

## Installation

- Clone repository ini dengan :

    ```sh
    git clone https://github.com/Giafn/Depublic.git
    cd Depublic
    ```

- Copy .env.example dan ubah namanya menjadi .env :

    ```sh
    cp .env.example .env
    ```
- konfigurasi .env sesuai kredensial Postgresql dan redisnya, untuk jwt dan encrypt bisa di biarkan default atau di ubah (*ENCRYPT_SECRET_KEY, dan ENCRYPT_IV HARUS 16 karakter)

- Download package golang yang dibutuhkan dengan :

    ```sh
    go mod tidy
    ```
- Jalankan Postgresql, dan Redis server

- Jalankan migrasi dengan make tool :

    ```sh
    make migration-up
    ```
- Jalankan server echo untuk memulai develop
    - bisa dengan make command
    ```sh
    make run-server
    ```

    - atau dengan
    ```sh
    go run cmd/app/main.go
    ```