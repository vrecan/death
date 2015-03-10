# death
<p>Simple library to make it easier to manage the death of your application.
Example code:
</p>
```
import (
        DEATH "github.com/vrecan/death"
        SYS "syscall"
)
death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM) //pass the signals you want to end your application
defer death.Close()

//when you want to block for shutdown signals
death.WaitForDeath() // this will finish when a signal of your type is sent to your application
```

## Close other structs when shutdown has been signaled
<p>One simple feature of death is that it can also close other objects when shutdown starts</p>
```
import (
        DEATH "github.com/vrecan/death"
        SYS "syscall"
)
death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM) //pass the signals you want to end your application
defer death.Close()

objects := make([]DEATH.Closable,0)

objects = append(objects, &{}newType) // this will work as long as the type implements a Close method

//when you want to block for shutdown signals
death.WaitForDeath(objects) // this will finish when a signal of your type is sent to your application
```
