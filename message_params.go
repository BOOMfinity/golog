package golog

import "github.com/BOOMfinity/go-utils/gpool"

type Params []*[2]any

var paramsPool = gpool.New[[2]any]()
