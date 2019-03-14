# leaderboard-app - a serverless application

[![OpenFaaS](https://img.shields.io/badge/openfaas-cloud-blue.svg)](https://www.openfaas.com)

This application is an example of how to write a Single Page App (SPA) with a Serverless approach. It provides a live leaderboard for your GitHub organisation or repos showing comments made and issues opened by your community and contributors.

* The front-end is written with Vue.js
* The backing data-store data is Postgres with a remote DBaaS or in-cluster deployment

See a live example tracking the openfaas/openfaas-incubator organizations: [https://alexellis.o6s.io/leaderboard-page](https://alexellis.o6s.io/leaderboard-page)

To test out the functionality comment on this issue: [Issue: Let's test the leaderboard!](https://github.com/openfaas/org-tester/issues/18)

Here's a preview of the app when the dark theme is enabled: 

![Dark Leaderboard example](docs/leaderboard-dark.png)

Subscribe to events by adding a webhook to the github-sub function:

![Subscribe](docs/subscribe.png)

## Functions

* github-sub

Receives webhooks from GitHub via an organization or repo subscription. Secured with HMAC by Alex Ellis

* leaderboard

Retrieves the current leaderboard in JSON by Alex Ellis

* leaderboard-page

Renders the leaderboard itself as a Vue.js app by Ken Fukuyama

## Schema

* [schema-1.0.sql](sql/schema-1.0.sql)

## Running locally

* Deploy OpenFaaS

* Grab the node8-express template

```
faas-cli template store pull node8-express
```

* Create the required secrets

```

export PASS=""
export USER=""
export HOST=""
export WEBHOOK=""   # As set on the webhook page on GitHub

faas-cli secret create leaderboard-app-secrets \
  --literal password="$PASS" \
  --literal username="$USER" \
  --literal host="$HOST" \
  --literal webhook-secret="$WEBHOOK"
```

* Deploy the stack.yml file with a prefix

Either edit the stack.yml and add the prefix for each function or run:

```
leaderboard => alexellis-leaderboard
github-sub => alexellis-github-sub
leaderboard-page => alexellis-leaderboard-page
```

* Deploy `of-router`:

Do this with auth turned off

https://github.com/openfaas/openfaas-cloud/tree/master/router

* Create entries in: `/etc/hosts`

```
alexellis.local-o6s.io    127.0.0.1
```

## Contributing & license

Please feel free to fork and star this repo and use it as a template for your own applications. The license is MIT.

To contribute see [CONTRIBUTING.md](./CONTRIBUTING.md)


