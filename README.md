# raptor
A Microservice Toolkit Developed Based On Golang.

```

                                                __----~~~~~~~~~~~----
                                     .  .   ~~//====......          _
                     -.            \_|//     |||\  ~~~~~~::::... /~
                  ___-==_       _-~o~  \/    |||  \          _/~~-
          __---~~~.==~||\=_    -_--~/_-~|-   |\   \        _/~
      _-~~     .=~    |  \-_    '-~7  /-   /  ||    \      /
    .~       .~       |   \ -_    /  /-   /   ||      \   /
   /  ____  /         |     \ ~-_/  /|- _/   .||       \ /
   |~~    ~~|--~~~~--_ \     ~==-/   | \~--===~~        .
            '         ~-|      /|    |-~\~~       __--~~
                        |-~~-_/ |    |   ~\_   _-~            /
                             /  \     \__   \/~                \__
                         _--~ _/ | .-~~____--~-/                  ~~=
                        ((->/~   '.|||' -_|    ~~-/ ,              . 
                                   -_     ~\      ~~---l__i__i__i--~~
                                   _-~-__   ~)  \--______________--~~
                                 //.-~~~-~_--~- |-------~~~~~~~~
                                        //.-~~~--    
```

## Quick Start

```
package main

import (
    "github.com/EAHITechnology/raptor/server"
)

func main(){
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)

    s, err := server.NewServer(ctx, "default", "./conf/server_config.toml")
    if err != nil {
    	fmt.Println("NewServer err:", err.Error())
    	return
    }
    
    if err := s.Run(ctx, cancel); err != nil {
    	fmt.Println("NewServer err:", err.Error())
    	return
    }
}
```

server will create a new service management manager. Create individual components through configuration.
When we need to use components, we can refer to the singleton object instantiated by the service manager in the corresponding package.If you use a service manager to build your application, then when using components you only need to:

```
import (
	"github.com/EAHITechnology/raptor/emysql"
)

func InsertDemoTask(d DemoModel) (int64, error) {
    // `cage-mysql` is the mysql name in the configuration file.

	client, err := emysql.GetClient("cage-mysql")
	if err != nil {
		return 0, err
	}

	if err := client.GetMaster().Table("demo_task").Create(&d).Error; err != nil {
		return 0, err
	}

	return d.Id, nil
}
```

No need to worry about initialization anymore.

## organizational structure
The organizational form of raptor is very simple, and functions are divided according to the package dimension, following the principle of single responsibility.You can find their interface initialization functions under each package

```
.
├── LICENSE
│
├── README.md
│
├── balancer            # load balancer
│
├── breaker             # fuse
│
├── config              # Parse the configuration file
│
├── context_trace       # traceId
│
├── distributed_lock    # Distributed locks implemented in multiple ways
│
├── double_buffer       # double pointer queue
│
├── elog                # log
│
├── emq                 # mq client
│
├── emysql              # mysql client
│
├── enet                # 
│
├── eredis              # redis client
│
├── erpc                # rpc client
│
├── example_config
│
├── limiter             # current limiter
│
├── server              # Service Governance 
│
├── service_discovery   
│
├── taskflow            # Simple workflow determinator, check DAG and topological sort
│
└── utils
```

You can also instantiate the component yourself, like this:

```
import (
    "github.com/EAHITechnology/raptor/elog"
)

func NewLogger(){
    lc:=&elog.LogConfig{}
    logger, err := elog.NewLogger(lc)
    if err != nil {
        ....
    }

    logger.Debug("xxxxx")
}
```