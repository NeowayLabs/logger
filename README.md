# Logger

[![Build Status](https://travis-ci.org/NeowayLabs/logger.svg?branch=master)](https://travis-ci.org/NeowayLabs/logger)

This package can help you add some log to your application. We have four different levels of log, **Debug**, **Info**,
**Warn** and **Error**, by default Debug will be *discarded*, Info and Warn will be redirect to *Stdout* and Error will
be redirect to *Stderr*. **Note** Error Level will always be redirect to Stderr, and you cannot disable that.

This module have a default logger instance with empty Namespace to make easy you use it without any additional line,
like we show below
```
package main

import log "github.com/NeowayLabs/logger"

func main() {
    log.Debug("number=%d string=%s...", 12, "test debug") // discarded
    log.Info("number=%d string=%s...", 10, "test info")   // [INFO] number=10 string=test info...
    log.Warn("number=%d string=%s...", 8, "test warn")    // [WARN] number=8 string=test warn...
    log.Error("number=%d string=%s...", 6, "test error")  // [ERROR] number=6 string=test error...
    // Fatal will abort program
    log.Fatal("number=%d string=%s...", 4, "test fatal")  // [FATAL] number=4 string=test error...
}
```

You can choose which level will be discarded or what will be shown calling ```SetLevel()``` passing
```logger.LevelDebug```, ```logger.LevelInfo```, ```logger.LevelWarn``` or ```logger.LevelError```. You can create new
instances with namespace if you want, to get new one call ```logger.Namespace("NAMESPACE)```.

You can use environment variable to set level instead call ```SetLevel``` manually, export ```SEVERINO_LOGGER``` with
```debug```, ```info```, ```warn``` and ```error```, this variable will set level to default namespace logger. To set
only of specifc module you can export ```SEVERINO_LOGGER_MY_MODULE```, if you don't do that, the level of default will
be used.
**NOTE:** the module name will be replace "-" and "." to "\_" and will be uppercase. If your module is: "vendor.my-module"
your environment variable will be "SEVERINO_LOGGER_VENDOR_MY_MODULE"

Take a look at following examples:

```
package main

import "github.com/NeowayLabs/logger"

func main() {
    // using default logger instance
    logger.Debug("number=%d string=%s...", 12, "test debug") // default debug is discarded
    logger.Warn("number=%d string=%s...", 8, "test warn")    // [WARN] number=8 string=test warn...

    // set level to show debug
    logger.SetLevel(logger.LevelDebug)
    logger.Debug("number=%d string=%s...", 12, "test debug") // [DEBUG] number=12 string=test debug...
    logger.Warn("number=%d string=%s...", 8, "test warn")    // [WARN] number=8 string=test warn...


    // getting new instance with namespace
    log := logger.Namespace("my-module")
    log.Debug("number=%d string=%s...", 12, "test debug") // default debug is discarded
    log.Warn("number=%d string=%s...", 8, "test warn")    // <my-module> [WARN] number=8 string=test warn...

    // set level to show debug
    log.SetLevel(logger.LevelDebug)
    log.Debug("number=%d string=%s...", 12, "test debug") // <my-module> [DEBUG] number=12 string=test debug...
    log.Warn("number=%d string=%s...", 8, "test warn")    // <my-module> [WARN] number=8 string=test warn...
}
```

### Use new handlers

You can create new handler to log in the different ways, you can implement log handler to send any kind of
alert, like send an email when error occour, or write warns to file. To do that you need implement at least of following
interface.
An example of our default implementation you found [here](http://github.com/NeowayLabs/logger/blob/master/handler.go)

* [Debug Interface](http://github.com/NeowayLabs/logger/blob/master/logger.go#L33)
* [Info Interface](http://github.com/NeowayLabs/logger/blob/master/logger.go#L37)
* [Warn Interface](http://github.com/NeowayLabs/logger/blob/master/logger.go#L41)
* [Error Interface](http://github.com/NeowayLabs/logger/blob/master/logger.go#L45)
* [Fatal Interface](http://github.com/NeowayLabs/logger/blob/master/logger.go#L49)
* [Init Interface](http://github.com/NeowayLabs/logger/blob/master/logger.go#L29) this function will be called
when you add your handler to logger instance and always ```setLevel``` was called


### HTTP handler

To avoid you have to restart your app to change level of your logger, we develop a HTTP Handler to you control all
the levels of your app, we have an usage example:

```
package main

import (
    "net/http"

    "github.com/NeowayLabs/logger"
)

func main() {
    log := logger.Namespace("TEST")

    http.Handle("/logger/", logger.HTTPHandler()) // handle /logger/NAMESPACE
    http.Handle("/logger", logger.HTTPHandler()) // handle /logger
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
```

* You can get the list of all namespaces and its levels (key _default_ is the empty namespace, default one)
```
curl http://localhost:8080/logger
```

* You can get the level of a specific namespace
```
curl http://localhost:8080/logger/test
```

* You can set the level of a specific namespace
```
curl -XPUT http://localhost:8080/logger --data '{"namespace": "", "level": "debug"}' # set level of default namespace (empty)
curl -XPUT http://localhost:8080/logger --data '{"namespace": "test", "level": "debug"}' # set level of test namespace
curl -XPUT http://localhost:8080/logger/test --data '{"level": "warn"}' # set level of test namespace
```

* You can set the level of all namespaces
```
curl -XPUT http://localhost:8080/logger/all --data '{"level": "debug"}' # set level of all namespaces
```
