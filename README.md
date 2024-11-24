# Performance

```
BenchmarkColorSimple/WithoutModule-12        	          624957	      1915 ns/op	       0 B/op	       0 allocs/op
BenchmarkColorSimple/WithModule-12           	          620372	      1946 ns/op	       0 B/op	       0 allocs/op
BenchmarkColorSimpleThreaded/WithoutModule-12         	  582339	      2060 ns/op	       0 B/op	       0 allocs/op
BenchmarkColorSimpleThreaded/WithModule-12            	  585877	      2068 ns/op	       0 B/op	       0 allocs/op
BenchmarkColorParams2-12                              	  515204	      2343 ns/op	       0 B/op	       0 allocs/op
BenchmarkColorParams2Threaded-12                      	  540090	      2164 ns/op	       0 B/op	       0 allocs/op
BenchmarkColorParams10-12                             	  326755	      3607 ns/op	       0 B/op	       0 allocs/op
BenchmarkColorParams10Threaded-12                     	  514057	      2338 ns/op	       0 B/op	       0 allocs/op
BenchmarkColorWithFormatString-12                     	  607726	      2077 ns/op	      16 B/op	       1 allocs/op
BenchmarkColorWithFormatStringThreaded-12             	  582708	      2070 ns/op	      16 B/op	       1 allocs/op
BenchmarkColorWithFormatInt-12                        	  583782	      2096 ns/op	      16 B/op	       1 allocs/op
BenchmarkColorWithFormatIntThreaded-12                	  572977	      2124 ns/op	      16 B/op	       1 allocs/op

BenchmarkJsonSimple/WithoutModule-12                  	  459572	      2760 ns/op	       0 B/op	       0 allocs/op
BenchmarkJsonSimple/WithModule-12                     	  466590	      2699 ns/op	       0 B/op	       0 allocs/op
BenchmarkJsonSimpleThreaded/WithoutModule-12          	  531388	      2459 ns/op	       0 B/op	       0 allocs/op
BenchmarkJsonSimpleThreaded/WithModule-12             	  519331	      2421 ns/op	       0 B/op	       0 allocs/op
BenchmarkJsonParams2-12                               	  364587	      3316 ns/op	       0 B/op	       0 allocs/op
BenchmarkJsonParams2Threaded-12                       	  505298	      2347 ns/op	       0 B/op	       0 allocs/op
BenchmarkJsonParams10-12                              	  273933	      4322 ns/op	       0 B/op	       0 allocs/op
BenchmarkJsonParams10Threaded-12                      	  482919	      2542 ns/op	       0 B/op	       0 allocs/op
BenchmarkJsonWithFormatString-12                      	  382524	      2730 ns/op	      16 B/op	       1 allocs/op
BenchmarkJsonWithFormatStringThreaded-12              	  525858	      2268 ns/op	      16 B/op	       1 allocs/op
BenchmarkJsonWithFormatInt-12                         	  444204	      2786 ns/op	      16 B/op	       1 allocs/op
BenchmarkJsonWithFormatIntThreaded-12                 	  533432	      2295 ns/op	      16 B/op	       1 allocs/op
```