[![Build Status](https://travis-ci.org/czerwe/gobravia.svg?branch=master)](https://travis-ci.org/czerwe/gobravia)

# gobravia
Library to controll an Bravia (2014)


## Parameters

### gobravia.GetBravia
| name | description | 
|---|---|
| IP | IP of Bravia | 
| PIN | PIN for remote configuration | 
| MAC | MAC address of host  | 

### gobravia.Poweron
| name | description | 
|---|---|
| Broadcast | Broadcast of network | 


## commands

The commands stored in an map, the key is an alias (e.g. up, down, num1).

```json
{
  "<alias>": "<command>"
}
```
The commands can be read by this
```go
bravia = gobravia.GetBravia("10.0.0.11", "0000", "FC:F1:52:F6:FA:8F")
bravia.GetCommands()
```
and can be printed to stdout with
```go
bravia.PrintCodes()
```
### Send Commands

The command can be send either by the command or by the alias.

```go
bravia.SendCommand(command)
```

```go
bravia.SendAlias(alias)
```

## example

The commands stored in an map, the key is an alias (e.g. up, down, num1).

```go
package main

import (
  "github.com/czerwe/gobravia"
)

func main() {
  bravia = gobravia.GetBravia("10.0.0.11", "0000", "FC:F1:52:F6:FA:8F")
	bravia.GetCommands()
  
  bravia.Poweron("10.0.0.255")

  bravia.SendCommand(command)
  
}
```
