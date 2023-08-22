```shell
wrk -t1 -d10s -c200 -s ./wrk/signup.lua http://localhost:8080/users/signup
```