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

## Installation

```bash
# Clone the repo
git clone https://github.com/kmilanbanda/chirpy.git
```

## Usage

To use the program, first navigate to the directory of your cloned repo. Next, make sure the program is built 
using ```go build -o out```. Then, begin running the program with ```./out```. Now you may access the site
in a browser of your choosing at localhost:8080. However, it is recommended to use curl to add a body and 
headers to the requests that you will be making. 

### Examples
#### Register User
curl -X POST -H "Content-Type: application/json" -d '{"email": "fake.email@example.com", "password": "p4ssword"}' http://localhost:8080/api/users
#### Login
curl -X POST -H "Content-Type: application/json" -d '{"email": "fake.email@example.com", "password": "p4ssword"}' http://localhost:8080/api/login
#### Post Chirp
curl -X POST -H "Content-Type: application/json" -d '{"body": "I love to chirp"}' http://localhost:8080/api/chirps

### Full Endpoint Documentation
To be filled out later...
