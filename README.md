# xml_parser_comparison



```s
Running tool: /usr/local/bin/go test -benchmem -run=^$ -coverprofile=/var/folders/z5/cj4zjtv15wn8yt53qzh57lnr0000gp/T/vscode-goj9wKIi/go-code-cover -bench . xml_parser_comparison

goos: darwin
goarch: amd64
pkg: xml_parser_comparison
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkStringBased-12         	  141146	      7661 ns/op	   12352 B/op	       5 allocs/op
BenchmarkEtreeBased-12          	    5790	    206467 ns/op	   93800 B/op	    1323 allocs/op
BenchmarkXmlEncodingBased-12    	    6192	    195523 ns/op	   75089 B/op	    1105 allocs/op
PASS
coverage: 78.7% of statements
ok  	xml_parser_comparison	3.849s
```