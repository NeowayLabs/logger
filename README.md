# Logger

This package can help you add some log to your application. We have four different levels of log, **Debug**, **Info**,
**Warn** and **Error**, by default Debug will be *discarded*, Info and Warn will be redirect to *Stdout* and Error will
be redirect to *Stderr*. **Note** Error Level will always be redirect to Stderr, and you can disable that.

This module have a default logger instance with empty Namespace to make easy you use it without any additional line, like we
show below
```
package main

import log "gitlab.neoway.com.br/teahupoo/severino/logger"

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
**NOTE:** the module name will be replace "-" and "." to "_" and will be uppercase. If your module is: "vendor.my-module"
your environment variable will be "SEVERINO_LOGGER_VENDOR_MY_MODULE"

Take a look at following examples:

```
package main

import "gitlab.neoway.com.br/teahupoo/severino/logger"

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
    log.Warn("number=%d string=%s...", 8, "test warn")    // [WARN] number=8 string=test warn...

    // set level to show debug
    log.SetLevel(logger.LevelDebug)
    log.Debug("number=%d string=%s...", 12, "test debug") // [DEBUG] number=12 string=test debug...
    log.Warn("number=%d string=%s...", 8, "test warn")    // [WARN] number=8 string=test warn...
}
```
