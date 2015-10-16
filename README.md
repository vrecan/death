# Death [![Build Status](https://travis-ci.org/vrecan/death.svg?branch=master)](https://travis-ci.org/vrecan/death)

[![Join the chat at https://gitter.im/vrecan/death](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/vrecan/death?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
<p>Simple library to make it easier to manage the death of your application.</p>

## Get The Library

```bash
go get github.com/vrecan/death
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
)

func main() {
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM) //pass the signals you want to end your application
	objects := make([]DEATH.Closable, 0)

	objects = append(objects, &NewType{}) // this will work as long as the type implements a Close method

	//when you want to block for shutdown signals
	death.WaitForDeath(objects...) // this will finish when a signal of your type is sent to your application
}

type NewType struct {
}

func (c *CloseMe) Close() error {
	return nil
}
```
