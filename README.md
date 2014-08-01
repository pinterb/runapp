runapp
======

Run your applications...with confidence(?)

### The idea
I really like the idea of [The Twelve-Factor App](http://12factor.net/). So the idea here is to enforce consistent application configuration 
by creating a wrapper command line script.  The approach this wrapper takes is fairly simplistic:

* Parse the command line arguments where the first argument is expected to be the application or command that you want to run (e.g "java")
* The parsed command line arguments become the initial flags that ultimately get passed to the application
* If one of the command line arguments is *--env_prefix*, that is the application prefix used to identify any environment variables that should also be parsed.
* If the *--env_prefix* flag is set, then parse any environment variables starting with that prefix.  By default, DO NOT override any command line arguments that were already set. ***Note:***
that the *--env_prefix* flag gets dropped...and does NOT get passed into the application you are trying to execute.
* Finally, if a .env file exists, parse its contents as well.  Again, by default, DO NOT override any flags that may have already been set by either command line or environment variable.  ***Important:*** Currently, the only file format that is supported is the "ini" file type.  Sections are supported by prepending the section name to the variable name (separated by an underscore).

### Status
* Very much an alpha-quality, [minimum viable product (MVP)](http://en.wikipedia.org/wiki/Minimum_viable_product).  By design, it's very opinionated.  But it has its warts as well. 
* This is my very first attempt at writing [Go](http://golang.org/). So my code is definitely a hack.  But the built-in tooling (e.g. go fmt, lint, vet) is awesome!
* Not a lot of testing has been done.  And no automated tests have been created.
* There is a third-party dependency. But no vendoring has been incorporated in the build.  You must go get the dependencies.


### Example

```bash
make
```
```java
bin/runapp java -jar tests/java/bin/hello.jar --env_prefix=MY_APP_ --db_user=bart
```

### To-Do's
* Too many to mention really.  But this approach isn't for everyone.  For example, any environment variables or config file variables parsed will get passed into the application for execution.
This is a bit simplistic and there may be ways to pre-configure your application to run (e.g create a .runapp config file with supported flags) more intelligently.
* Give some thought to creating sub-commands around known languages (e.g. java, ruby, python).  For example, this might be useful for java as something like *-XX:MaxPermSize* is currently not handled as you'd expect (you need to specify XX:MaxPermSize as your variable name).  In fact, these type of java options may only work from the command line as all environment and config file variables are created as long flags (e.g "--").
* With Docker containers becoming popular, more and more configuration will be available via solutions like [etcd](https://github.com/coreos/etcd) or perhaps [consul](http://www.consul.io/).  This application should identify a way to support these types of configuation sources.
* Tests and vendoring.
