# Test task from Synergy

```
    docker build -t synergy_test_task .
    docker run -dp 8080:8080 synergy_test_task
    make go-build-client
    make client
```
