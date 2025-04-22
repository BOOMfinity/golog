## Configurable variables

```
stackTraceBufferSize:            1024,
loggerParametersSliceAllocation: 25,
messageParameterSliceAllocation: 25,
messageBufferSize:               1024,
engineBufferSize:                2048,
overrideMinimumMessageLevel:     0,
disableTerminalColors:           false,
marshalDetails:                  true,
detailsBufferSize:               256,
includeStackOnError:             false,
```

Use `golog.Config.SetXxx` methods to change the configuration. You should call them at the top of the main function.

## Performance

```
goos: windows
goarch: amd64
pkg: github.com/BOOMfinity/golog/v2
cpu: AMD Ryzen 5 4600H with Radeon Graphics         
BenchmarkColor
BenchmarkColor/JustMessage
BenchmarkColor/JustMessage-12         	 2549898	       471.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/Modules
BenchmarkColor/Modules/1
BenchmarkColor/Modules/1-12           	 2433614	       507.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/Modules/3
BenchmarkColor/Modules/3-12           	 2087161	       597.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/Modules/5
BenchmarkColor/Modules/5-12           	 1911147	       662.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/Modules/10
BenchmarkColor/Modules/10-12          	 1540208	       864.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/Modules/20
BenchmarkColor/Modules/20-12          	 1000000	      1093 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/WithDetails
BenchmarkColor/WithDetails/ParseToJSON=true
BenchmarkColor/WithDetails/ParseToJSON=true/string
BenchmarkColor/WithDetails/ParseToJSON=true/string-12  	 1926484	       627.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/WithDetails/ParseToJSON=true/int
BenchmarkColor/WithDetails/ParseToJSON=true/int-12     	 1985318	       603.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/WithDetails/ParseToJSON=true/float
BenchmarkColor/WithDetails/ParseToJSON=true/float-12   	 1612898	       750.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/WithDetails/ParseToJSON=true/slice
BenchmarkColor/WithDetails/ParseToJSON=true/slice-12   	 1591832	       745.3 ns/op	      24 B/op	       1 allocs/op
BenchmarkColor/WithDetails/ParseToJSON=false
BenchmarkColor/WithDetails/ParseToJSON=false/string
BenchmarkColor/WithDetails/ParseToJSON=false/string-12 	 2194882	       539.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/WithDetails/ParseToJSON=false/int
BenchmarkColor/WithDetails/ParseToJSON=false/int-12    	 2155323	       557.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/WithDetails/ParseToJSON=false/float
BenchmarkColor/WithDetails/ParseToJSON=false/float-12  	 1801533	       666.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/WithDetails/ParseToJSON=false/slice
BenchmarkColor/WithDetails/ParseToJSON=false/slice-12  	  995528	      1157 ns/op	      64 B/op	       6 allocs/op
BenchmarkColor/UserMessage
BenchmarkColor/UserMessage/10
BenchmarkColor/UserMessage/10-12                       	 2491965	       470.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/UserMessage/25
BenchmarkColor/UserMessage/25-12                       	 2465709	       485.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/UserMessage/50
BenchmarkColor/UserMessage/50-12                       	 2386390	       496.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/UserMessage/100
BenchmarkColor/UserMessage/100-12                      	 2244390	       528.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/UserMessage/200
BenchmarkColor/UserMessage/200-12                      	 2105282	       567.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/UserMessage/400
BenchmarkColor/UserMessage/400-12                      	 1810768	       652.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkColor/UserMessage/800
BenchmarkColor/UserMessage/800-12                      	 1066604	      1129 ns/op	     898 B/op	       1 allocs/op
BenchmarkColor/UserMessage/1600
BenchmarkColor/UserMessage/1600-12                     	  695284	      1782 ns/op	    1797 B/op	       1 allocs/op
BenchmarkJSON
BenchmarkJSON/JustMessage
BenchmarkJSON/JustMessage-12                           	 1366868	       808.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/Modules
BenchmarkJSON/Modules/1
BenchmarkJSON/Modules/1-12                             	 1296646	       969.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/Modules/3
BenchmarkJSON/Modules/3-12                             	 1000000	      1092 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/Modules/5
BenchmarkJSON/Modules/5-12                             	 1000000	      1259 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/Modules/10
BenchmarkJSON/Modules/10-12                            	  781396	      1748 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/Modules/20
BenchmarkJSON/Modules/20-12                            	  528822	      2670 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/WithDetails
BenchmarkJSON/WithDetails/string
BenchmarkJSON/WithDetails/string-12                    	 1370022	       879.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/WithDetails/int
BenchmarkJSON/WithDetails/int-12                       	 1379214	       873.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/WithDetails/float
BenchmarkJSON/WithDetails/float-12                     	 1000000	      1022 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/WithDetails/slice
BenchmarkJSON/WithDetails/slice-12                     	 1000000	      1020 ns/op	      24 B/op	       1 allocs/op
BenchmarkJSON/UserMessage
BenchmarkJSON/UserMessage/10
BenchmarkJSON/UserMessage/10-12                        	 1456159	       822.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/UserMessage/25
BenchmarkJSON/UserMessage/25-12                        	 1391900	       853.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/UserMessage/50
BenchmarkJSON/UserMessage/50-12                        	 1295318	       924.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/UserMessage/100
BenchmarkJSON/UserMessage/100-12                       	 1000000	      1015 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/UserMessage/200
BenchmarkJSON/UserMessage/200-12                       	 1000000	      1163 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/UserMessage/400
BenchmarkJSON/UserMessage/400-12                       	  789945	      1491 ns/op	       0 B/op	       0 allocs/op
BenchmarkJSON/UserMessage/800
BenchmarkJSON/UserMessage/800-12                       	  392461	      2628 ns/op	     900 B/op	       1 allocs/op
BenchmarkJSON/UserMessage/1600
BenchmarkJSON/UserMessage/1600-12                      	  281132	      4219 ns/op	    1800 B/op	       1 allocs/op
```