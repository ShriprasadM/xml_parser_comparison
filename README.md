# xml_parser_comparison



```s
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -coverprofile=/var/folders/_q/40r7ttwx7gvcgy_nv46hz0x40000gp/T/vscode-goQUY6HW/go-code-cover -bench . xml_parser_comparison

goos: darwin
goarch: amd64
pkg: xml_parser_comparison
cpu: Intel(R) Core(TM) i5-8259U CPU @ 2.30GHz
BenchmarkStringBased-8        	  504448	      2377 ns/op	    3072 B/op	       1 allocs/op
BenchmarkEtreeBased-8         	   10000	    115355 ns/op	   47504 B/op	     659 allocs/op
BenchmarkXmlEncodingBased-8   	   10000	    117195 ns/op	   37812 B/op	     589 allocs/op
PASS
coverage: 71.7% of statements
ok  	xml_parser_comparison	3.589s
```