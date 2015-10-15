# Scribble

[![Circle CI](https://circleci.com/gh/keichi/scribble.svg?style=svg)](https://circleci.com/gh/keichi/scribble)

Simple web application to take notes in Markdown.

## How to run Scribble

### Get dependencies

```
$ go get -v ./...
$ cd static
$ npm install
$ bower install
```

### Build & run (development)

```
$ go run main.go
$ cd static
$ grunt serve
```

### Build & run (production)

```
$ cd static
$ grunt build
$ cd ..
$ go run main.go
```

## Architecture

### Server side

- kami (Web Application Framework)
- context (Request global variable)
- gorp (ORM)
- testify (Assertion)
- goamz (AWS S3 client)

### Client side

- AngularJS (Client-side MVC)
- Bootstrap (CSS Framework)
- AngularUI (UI components for AngularJS)
- Ace Editor (Rich source code editor)
- hightlight.js (Source code hightlighter)
- marked (Markdown renderer)
- keymaster (Keyboard shortcuts)

