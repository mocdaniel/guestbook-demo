# Contributing

When providing bug fixes, ideas for improvement or something entirely different,
you might have to create a proper development environment locally. 

This document guides you through setting up a working local environment and points out 
some known caveats.

## Prerequisites

Software needed for this guide to work out of the box:

* golang
* npm
* a webbrowser of your choice
* Docker (optional)

## Initializing the Frontend Project

The frontend is written in **Typescript** using **VueJS** and **TailwindCSS**. After cloning the repository, 
navigate to `./app/frontend` and fetch all needed dependencies.

```
cd ./app/frontend
npm ci
```

Spin up the frontend afterwards and visit [http://localhost:5173](http://localhost:5173). You should see the (empty) Web UI of the project.

```
npm run dev

> guestbook-demo@1.0.0 dev
> vite


  VITE v3.1.8  ready in 238 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

## Initializing the Backend Project

The backend is written in **Go**. After cloning the repository, fetch all needed dependencies from the repository's root directory.

```
go get ./app
```

## Running the Frontend on Hot-Reload

During development, especially when introducing changes to the UI, you might want to see the impact of those changes immediately.

For this to work, you will need to run the frontend separately, using the provided `watch` command:

```
cd ./app/frontend
npm run watch
```

This command will constantly look for changes within the VueJS project, and (re-)build the static files on detected changes. The output can be found at `./app/frontend/dist` and will be served by our Go application automatically.

To check functionality of both, backend and frontend, you can then run `go run ./app` and visit [http://localhost:8080](http://localhost:8080).

Please note that even in developer mode, a working connection to both, **Redis** and **PostgreSQL** services needs to be provided. For more information, see the section about [Additional Services](#additional-services) below.

## Additional Services

For the application to work, a **PostgreSQL** database as well as a **Redis** cache need to be running even in developer environments. This can be done easily via Docker:

```
# Run a PostgreSQL instance matching the application's default settings
docker run -d -p 5432:5432 --name postgres -e POSTGRES_USER=guestbook -e POSTGRES_PASSWORD=password -e POSTGRES_DATABASE=guestbook postgres:12

# Run a Redis instance matching the application's default settings
docker run -d -p 6379:6379 --name redis redis
```

Once you're done, you can clean up the containers like this:

```
docker stop redis postgres
docker rm redis postgres
```