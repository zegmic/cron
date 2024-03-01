# Cron Expression Parser
The program takes a cron pattern and a program to run as an argument e.g. 
```bash
./bin/cron */15 0 1,15 * 1-5 /usr/bin/find
```
parses it and converts to an output
```bash
minute         0 15 30 45 
hour           0 
day of month   1 15 
month          1 2 3 4 5 6 7 8 9 10 11 12 
day of week    1 2 3 4 5 
command        /usr/bin/find
```

## Building
In order to build the project Go compiler (min 1.21) is required. You can download it from https://go.dev/dl/

Once Go is installed the following commands can be run in a project directory

Running a test suite
```bash
make test
```

Building the binary
```bash
make build
```

Once the project has compiled you can run it
```bash
./bin/cron */15 0 1,15 * 1-5 /usr/bin/find
```
