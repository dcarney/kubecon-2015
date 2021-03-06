# Kubernetes Is <br>For Monoliths Too
##### KubeCon 2015, Nov 10, 2015

***

Steve Reed (@_smreed)

---

# About Me

Early Kubernetes adopter:

* v0.5 in GCE ~Nov 2014
* [#3965](https://github.com/kubernetes/kubernetes/issues/3965) "GCE PD Data Loss"
* 12 commits, ~143 lines survive today
* One of the first production workloads

---

# Outline

* Its service discovery can eliminate a lot of configuration. From
  dev/test/staging/prod copies of everything, to port selection nightmares.
* Example: dev/test/staging/prod can be namespaces, and you can otherwise
  "hardcode" service names. Required dependent services are first class
  citizens in your infrastructure, this in my opinion isn't a bad thing to do.
* Even if you have a database you don't want to containerize, it can still
  be a service in Kubernetes, and utilize the same benefits.
* When your configuration absorbs too many details about your surrounding
  infrastructure, you can wrap that complexity up in an ambassador container.
  For example, wrap up smtp auth w/ an open relay. Run twemproxy along-side
  your app.
* Replace/add features using adapters (nginx for ssl termination)
* Show a sample monolithic application configuration in its simplest form.
  Show how Kubernetes still allows for flexibility, app can be run in a 
  dev or prod mode w/ no change. Walk through changes made to mattermost,
  for example.

---

# Motivation

![k8s-vs-microservices](./assets/k8s-vs-microservices.png)
<small>Google Searches For "Kubernetes" and "Microservices"</small>

* Everybody should be thinking of Kubernetes
* Microservices are not for everybody
* Monolithic applications can benefit
* Kubernetes is not (just) "for" microservices

<!--
* (?) Kubernetes Isn't just for microservices
* There are reasons to move your monolith into Kubernetes
-->

---

# Monoliths <br>These Days

* Application deployed as a single executable or package.
* An amalgam of orthogonal concerns.
* May still utilize many other external services (DB, cache, SMTP, etc)

<!--
* Briefly, What's a monolith? 
-->

---

# Helmsman <br>Of The Titanic

The metaphor doesn't entirely fit, but it sounds cool.

* `Services` elevate your external dependencies to 1st class status
* 
Can Kubernetes help with

* day to day development
* staging and testing

---

# Title

A monolithic application with many external dependencies can lead to complicated
development environments.

---

# Spaghetti.Conf (1/n)

Common to have "dev", "test", and "prod" configuration for the same application.

Kubernetes can provide the same flexibility with namespaces.

---

# Spaghetti.Conf (2/n)

```
# prod-app.properties
db.host = database.prod
db.port = 3306
```

---

# Sandbox Slide

<pre><code class="go hljs">package main

import "fmt"

func main() {
  fmt.Println("Hello!")
}</code></pre>

