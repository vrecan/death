# Death [![Build Status](https://travis-ci.org/vrecan/death.svg?branch=master)](https://travis-ci.org/vrecan/death) [![Coverage Status](https://coveralls.io/repos/github/vrecan/death/badge.svg?branch=master)](https://coveralls.io/github/vrecan/death?branch=master)

[![Join the chat at https://gitter.im/vrecan/death](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/vrecan/death?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
<p>Simple library to make it easier to manage the death of your application.</p>

## Get The Library

Use gopkg.in to import death based on your logger.

Version | Go Get URL | source | doc | Notes |
--------|------------|--------|-----|-------|
2.x     | [gopkg.in/vrecan/death.v2](https://gopkg.in/vrecan/death.v2)| [source]() | [doc]() | This supports loggers who _do not_ return an error from their `Error` and `Warn` functions like [logrus](https://github.com/sirupsen/logrus)
1.x     | [gopkg.in/vrecan/death.v1](https://gopkg.in/vrecan/death.v1)| [souce](https://github.com/vrecan/death/tree/v1.0) | [doc](https://godoc.org/gopkg.in/vrecan/death.v1) | This supports loggers who _do_ return an error from their `Error` and `Warn` functions like [seelog](https://github.com/cihub/seelog)



Example
```bash
go get gopkg.in/vrecan/death.v2
```
## Use The Library

```go
package main

import (
	DEATH "github.com/vrecan/death"
	SYS "syscall"
)

func main() {
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM) //pass the signals you want to end your application
	//when you want to block for shutdown signals
	death.WaitForDeath() // this will finish when a signal of your type is sent to your application
}
```

### Close Other Objects On Shutdown
<p>One simple feature of death is that it can also close other objects when shutdown starts</p>

```go
package main

import (
	DEATH "github.com/vrecan/death"
	SYS "syscall"
	"io"
)

func main() {
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM) //pass the signals you want to end your application
	objects := make([]io.Closer, 0)

	objects = append(objects, &NewType{}) // this will work as long as the type implements a Close method

	//when you want to block for shutdown signals
	death.WaitForDeath(objects...) // this will finish when a signal of your type is sent to your application
}

type NewType struct {
}

func (c *NewType) Close() error {
	return nil
}
```

### Or close using an anonymous function

```go
package main

import (
	DEATH "github.com/vrecan/death"
	SYS "syscall"
)

func main() {
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM) //pass the signals you want to end your application
	//when you want to block for shutdown signals
	death.WaitForDeathWithFunc(func(){ 
		//do whatever you want on death
	}) 
}
```
