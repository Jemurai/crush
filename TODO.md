# TODO Items for crush


```sh
mk@om crush % go run main.go examine --directory /Users/mk/area54/service --tag injection --ext .clj
```


mk@om crush % go run main.go examine --directory /Users/mk/area54/service --tag clojure > ./testing-clojure.json                

Fix stuff...

mk@om crush % go run main.go examine --directory /Users/mk/area54/service --tag clojure --compare ./testing-clojure.json

mk@om crush % ./build.sh github.com/jemurai/crush

go run main.go examine --debug true  --directory /Users/mk/area51/jasp/jasp-api --threshold 6.0

## Checks To Implement

- MD5
- SHA-1
- Argon
- AES
- Bcrypt

- Redis_pass

- jdbc/query + (str