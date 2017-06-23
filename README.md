# utmpsave
Dumps linux lastlog into a human readable file. 
Written in pure go, with no bash forking makes it efficient and blazing fast.
utmpsave can be scheduled. (won't hold duplicate records in dump file)

Build steps:

1. Clone the repository in your `$GOPATH/src/`
2. `go build utmpsave.go`
3. `./utmpsave`


```
Usage:
  utmpsave [flags]

Flags:
  -d, --destination string   Destination dump file path (default "/tmp/utmpdump.dmp")
  -h, --help                 help for utmpdump
  -s, --source string        Source wtmp file path (default "/var/log/wtmp")
  -v, --version              Prints version info
```
