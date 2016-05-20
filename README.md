# Enliven

## An easy-to-use web application framework in Go

### Features:

* Developed to deliver similar features to the "Django" framework (python)
* Wraps and provides gorilla/mux for routing
* Uses jinzhu/gorm for database management/interaction
* Contains a fork of qor/admin, an administration panel (similar to Django w/ Suit)
* Middleware management inspired by codegangsta/negroni
* Dependency Injection via context provided to handlers/middleware
* Session management with multiple storage drivers
* User account management, including user roles and permissions
* Static asset serving via true filesystem or embedded assets using Statik
* API for writing packaged enliven "apps" which can be easily added to any Enliven app


## Project State

* Enliven is currently in an alpha/experimental development stage.
* Once the framework is feature complete there will be a short beta to gather feedback and fix issues.
* Release date is "when it's done," but will adhere to "release early, release often."
