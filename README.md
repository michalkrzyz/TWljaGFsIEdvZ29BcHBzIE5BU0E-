# url-collector
STARTING APP
============


To start application locally go to app dir and use:


~$ go run .


To start docker application, build it and run in following steps:


~$ docker-compose build

...

~$ docker-compose up


Aplication has following configuration environment variables:


URL\_COLLECTOR\_PORT: PORT number for collector
NASA\_API\_KEY: API\_KEY for nasa api service (default DEMO\_KEY)
CONCURRENT\_REQUESTS: number of maximum requests served in parallel


UTs
===


To run UTs suite go to app directory and issue following command:


~$ go test -v ./...


INTEGRATION TESTS
=================


To run integration tests, go to integration\_test directory,
start your application (either in docker or locally)
and issue following command:


~$ robot test.robot


URL\_COLLECTOR\_PORT variable from test can be easily changed,
use following example:


~$ robot -v URL\_COLLECTOR\_PORT:8092 test.robot


Dependency of integration tests:
python3, python3-robotframework, python3-robotframework-requests
