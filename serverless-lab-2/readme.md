# README

## Buoi 2

### Mô tả

Viết 4 API, gọi vào database postgres, oracle or mysql

### Các API

#### API 1: Create User

Tạo user với username, name, phone. Phải check username tồn tại duy nhất trong table, không sử dụng unique của database

##### Input

```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "username": {{string}},
        "name": {{string}},
        "phone": {{string}}
    }
}
```

##### Output

```
{
    "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "responseCode": {{string}},
    "responseMessage": {{string}}
}
```

#### API 2: Update User

Update user by username. Thông tin update là name và phone.

##### Input

```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "username": {{string}},
        "name": {{string}},
        "phone": {{string}}
    }
}
```

##### Output

```
{
    "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "responseCode": {{string}},
    "responseMessage": {{string}},
}
```

#### API 3: Delete User

Xóa user by username

##### Input

```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "username": {{string}}
    }
}
```

##### Output

```
{
    "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "responseCode": {{string}},
    "responseMessage": {{string}},
}
```

#### API 4: Get User Detail

Lấy thông tin user by username

##### Output

```
{
    "responseId": {{uuid}},
    "responseTime": {{timeRPC3339}},
    "responseCode": {{string}},
    "responseMessage": {{string}},
    "data": {
        "username": {{string}},
        "name": {{string}},
        "phone": {{string}}
    }
}
```

### Yêu cầu

- Trong tất cả các API điều phải validate username, name, phone không được rỗng, KHÔNG được sài pointer
- Thiết kế API path cũng như method hợp lý.
- Script create table user

```
CREATE TABLE "users" (
    "id" bigserial,
    username character varying(50) COLLATE pg_catalog."default",
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    phone character varying(50) COLLATE pg_catalog."default",
    PRIMARY KEY ("id")
);
```


