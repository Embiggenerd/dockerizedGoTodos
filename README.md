# Dockerized go todos

A simple authenticated todos app written using only sql, the go standard lib, and docker.

## Visit App

Clone the repo, type docker-compose up. Uses port 80.

## In Detail

* Input some details to register. Only email and password are necessary. Check out error messages. 
* Login, hit the submit button. You can edit, delete todos
* Users are authenticated by a hex token, not user ID.
* On login, a new random hex value is created, and kept in a session table along with user ID
* To authenticate, middleware gets userid by looking up hex token in cookie
* Hex token is refreshed every login for security, old sessions erased for space
* This makes man in the middle attacks more difficult

## Demonstrates

* Ability to read documentation and source code.