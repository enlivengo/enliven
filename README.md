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

* Enliven is mainly a personal project custom tailored to my personal desires/needs.
* As such, I consider this a production-ready, mature project, and currently have some
* closed-source projects which are presently using it. I don't see myself writing a
* whole pile of documentation for it unless it magically gains a following and users
* request documentation.

* Future updates will consist of features I want to add to the framework in addition
* to bug fixes, etc. Again, I'm totally open to adding requested features as well,
* though that would require the aforementioned magical following gains.
