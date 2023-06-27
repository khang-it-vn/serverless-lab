
# Serverless-lab

## Golang Courses

### Buoi 1

#### Deploy Lambda with Golang, expose 2 REST APIs

All APIs have 2 default fields: `requestId`, `requestTime`

##### API-1

Input has 2 fields: `value1`, `value2`

```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "value1": {{number}},
        "value2": {{number}}
    }
}
```

Output: returns the sum of `value1` and `value2`

```
{
   "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "data": {
        "sum": {{value1+value2}}
    }
}
```

##### API-2

Input has 2 fields: `plaintText`, `secretKey`

```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "plaintText": {{string}},
        "secretKey": {{string}}
    }
}
```

Output: returns 1 field: `signature` using sha256 or hmacsha256 algorithm

```
{
   "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "data": {
        "signature": {{string}}
    }
}
```

##### API-3

Uses base64, input has 2 fields: `needEncode`, `needDecode`

```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "needEncode": {{string}},
        "needDecode": {{string}}
    }
}
```

Output: returns 2 fields: `outEncode` is the output of base64 field `needEncode`, `outDecode` is the output of field `needDecode`

```
{
   "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "data": {
        "outEncode": {{string}},
        "outDecode": {{string}}
    }
}
```


