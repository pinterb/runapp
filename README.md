runapp
======

Run your applications...with confidence(?)

### The idea
I really like the idea of [The Twelve-Factor App](http://12factor.net/). So the idea here is to enforce consistent application configuration 
by creating a wrapper command line script.  The approach this wrapper takes is fairly simplistic:

* Parse the command line arguments where the first argument is expected to be the application or command that you want to run (e.g "java")
* The parsed command line arguments become the initial flags that ultimately get passed to the application
* If one of the command line arguments is *--env_prefix*, that is the application prefix used to identify any environment variables that should also be parsed.
* If the *--env_prefix* flag is set, then parse any environment variables starting with that prefix.  By default, DO NOT override any command line arguments that were already set.
* Finally, if a .env file exists, parse its contents as well.  Again, by default, DO NOT override any flags that may have already been set by either command line or environment variable.  ***Important*** Currently, the only file format that is supported is the "ini" file type.  Sections are supported by prepending the section name to the variable name (separated by an underscore).

### Status
* Very much an alpha-quality, [minimum viable product (MVP)](http://en.wikipedia.org/wiki/Minimum_viable_product).  By design, it's very opinionated.  But it has its warts as well. 
* This is my very first attempt at writing [Go](http://golang.org/). So my code is definitely a hack.  But the built-in tooling (go fmt, lint, vet) is awesome!
* Not a lot of testing has been done.  And no automated tests have been created.
* There is a third-party dependency. But no vendoring has been incorporated in the build.  You must go get the dependencies.


### Example

```java
bin/runapp java -jar tests/java/bin/hello.jar --env_prefix=MY_APP_ --db_user=bart
```
