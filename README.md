Regex-Golang
====

Regex-Golang is an online Go regular expression tester

Inspired by [rubular](http://rubular.com/) and [regoio](https://regoio.herokuapp.com/).

It's currently deployed on App Engine at [https://regex-golang.appspot.com/](http://regex-golang.appspot.com/)

## Enhancements over regoio
* Do not write regex to logs
* Serve app assets from App Engine's static file handler (leverages Google's CDN)
* Upgrade to Bootstrap 3 and pull from CDN
* Log additional errors

## TODO

* Sharing (permalink)
* Add developer documentation
