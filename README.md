# Dockerized go todos

A simple authenticated todos app written using only sql, the go standard lib, and docker.

## Visit App

Clone the repo, type docker-compose up. Uses port 8080. Metrics can be viewed at /debug/vars, with 'requests' being a 
request counter. Visit the app and check that the number has changed.

## In Detail

### User Experience
* Input some details to register. Only email and password are necessary. Check out error messages. 
* Login, hit the submit button. You can edit, delete todos.

### User Authentication
* Users are authenticated by a hex token, not user ID.
* A session table has a userID and hex columns.
* On login, a new random hex value is created, inserted into a new session, and old session is deleted.
* To authenticate, middleware gets userid by looking up hex token in cookie.
* Hex token is refreshed every login for security, old sessions erased for space
* This makes man in the middle attacks more difficult

### Logging
* Middleware logs request data to json file as well as stdout
* Logs file is human readable, with newest at the top

### Metrics
* All standard pprof endpoints are available
* /debug/vars will also show number of requests through middleware

## Demonstrates

* Docker, docker-compose
* Basic sql, crud functionality
* User registration/login
* Metrics and observability