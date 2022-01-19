# xml_parser_comparison



```s
Running tool: /usr/local/bin/go test -benchmem -run=^$ -coverprofile=/var/folders/z5/cj4zjtv15wn8yt53qzh57lnr0000gp/T/vscode-goKLpmm1/go-code-cover -bench . xml_parser_comparison

goos: darwin
goarch: amd64
pkg: xml_parser_comparison
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkStringBased-12         	  141571	      8242 ns/op	   12982 B/op	      12 allocs/op
BenchmarkEtreeBased-12          	    5586	    234186 ns/op	   94560 B/op	    1324 allocs/op
BenchmarkXmlEncodingBased-12    	    5816	    205643 ns/op	   75370 B/op	    1109 allocs/op
PASS
coverage: 83.3% of statements
ok  	xml_parser_comparison	5.309s
```