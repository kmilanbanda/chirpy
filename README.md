# Chirpy

Chirpy is a Twitter-like HTTP server program with endpoints mirroring the basic features of Twitter/X.

## Features

 - Login, post chirps, get chirps etc.
 - Uses access, refresh, and api authentication tokens
 - Uses sql to handle the chirp/user data

## Motivation

This project was meant to help me become more familiar with HTTP connections from the server side.
I have a much better understanding now of how a server receives a request and transforms that body
into action.

## Requirements

 - Go
 - Postgres
 - Goose or similar program for migrating your database

## Installation

```bash
# Clone the repo
git clone https://github.com/kmilanbanda/chirpy.git
```

###Create a .env file and fill out the following variables:
 - DB\_URL="yourDatabaseConnectionString"
 - PLATFORM="dev"
 - MAX\_CHIRP\_LENGTH=140
 - SECRET="aLongGeneratedStringWithRandomCharacters" (mine was 88 characters long)
 - POLKA\_KEY="B26AE507C12A64AA4E78A7683E18371F" (a 32 bit hexadecimal string, try numbergenerator.org)

## Usage

To use the program, first navigate to the directory of your cloned repo. Make sure you add a ".env" file 
and fill out the variables in the Installation section. Use Goose and Postgres to migrate up your database.
Next, make sure the program is built using ```go build -o out```. Then, begin running the program with 
```./out```. Now you may access the site in a browser of your choosing at localhost:8080. 

###You can test the full functionality of this project using curl commands:

### Examples
#### Register User
curl -X POST -H "Content-Type: application/json" -d '{"email": "fake.email@example.com", "password": "p4ssword"}' http://localhost:8080/api/users
#### Login
curl -X POST -H "Content-Type: application/json" -d '{"email": "fake.email@example.com", "password": "p4ssword"}' http://localhost:8080/api/login
#### Post Chirp
curl -X POST -H "Content-Type: application/json" -d '{"body": "I love to chirp"}' http://localhost:8080/api/chirps

### Full Endpoint Documentation
    GET /api/healthz - returns server status
	GET /admin/metrics - get hits on the site
    POST /api/users - Creates user
	POST /api/login - login
	POST /admin/reset - resets databases
	POST /api/chirps - posts chirp
    GET /api/chirps - gets ALL chirps
	GET /api/chirps/{chirpID} - gets a specific chirp with {chirpID}
    POST /api/refresh - gets a new access token using a refresh token
    POST /api/revoke - revokes a refresh token
	PUT /api/users - updates a user's email and/or password
    DELETE /api/chirps/{chirpID} - deletes a chirp with {chirpID}
    POST /api/polka/webhooks" - allows a "third party" to upgrade a user to Chirpy Red

WIP: Endpoints will be further described with their appropriate request bodies at a later time


