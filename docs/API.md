# API 文档

以下接口为已经过测试的接口，可查看GraphQL Schema文件夹 [graph/schema](graph/schema) 查看全部接口

[TOC]

## Account

### 登录

Query

```graphql
mutation ($input: LoginInput!) {
    login(input: $input) {
        id
        token
        permission
        expire
    }
}
```

Input

```json
{
  "input": {
    "username": "<用户名>",
    "password": "<密码>"
  }
}
```

### 注册

Query

```graphql
mutation ($input: RegisterInput!) {
    register(input: $input)
}
```

Input

```json
{
  "input": {
    "username": "用户名",
    "email": "mail@example.com",
    "email_verify_code": "邮箱验证码",
    "password": "密码"
  }
}
```

### 发送邮件

Query

```graphql
mutation($input:SendEmailCodeInput!) {
    sendEmailCode(input:$input)
}
```

Input

```json
{
  "input": {
    "mail": "mail@example.com",
    "register": false
  }
}
```

### 找回密码

Query

```graphql
mutation ($input:ForgetInput!){
    forget(input:$input)
}
```

Input

```json
{
  "input": {
    "email": "mail@example.com",
    "password": "<新密码>",
    "email_verify_code": "<邮箱验证码>"
  }
}
```

## Memory

### 创建 Memory

Query

```graphql
mutation ($input:AddMemoryInput!){
    addMemory(input:$input)
}
```

Input

```json
{
  "input": {
    "title": "<标题>",
    "content": "<正文>"
  }
}
```

## 用户

### 当前用户

Query

```graphql
query{
    currentUser{
        id
        username
        avatar
        mail
        login_time
        create_time
        permission
        used
        subscribe{
            id
            name
            capacity
            available
            create_time
        }
    }
}
```

Input

```json
{
  "data": {
    "currentUser": {
      "id": "6447a3beb41134bfde45423b",
      "username": "<用户名>",
      "avatar": "[头像,默认为空字符串]",
      "mail": "mail@example.com",
      "login_time": 1682496790,
      "create_time": 1682416574,
      "permission": 0,
      "used": 0,
      "subscribe": {
        "id": "000000000000000000000000",
        "name": "免费版",
        "capacity": 104857600,
        "available": true,
        "create_time": 0
      }
    }
  }
}
```

## 资源相关

### 获取文件上传 Token

**有效时间** 60s  
Query

```graphql
mutation {
    # 文件名格式必须为 <16位或32位16进制字符串>(.png/.jpg/.jpeg/.gif/.webp或无后缀)
    # 推荐使用 32位 md5 散列值充当文件名
    getToken(fileName:"1234567890abcdef.png"){
        access_key
        secret_access_key
        session_token
        user_id
    }
    # 本地存储时仅需 session_token ，百度云对象存储需要使用到所有四个字段
}
```

Output

```json
{
  "data": {
    "getToken": {
      "access_key": "",
      "secret_access_key": "",
      "session_token": "11bc4458b9d9b85dfd278552ed762eea",
      "user_id": "6447a3beb41134bfde45423b"
    }
  }
}
```

### 上传本地存储文件

仅服务器配置 `storage_provider` 为 `local` 时可用  
**请求类型** `multipart/form-data`

#### 请求参数

_**operations**_

```json
{
  "query": "mutation($input:LocalUploadInput!){ localUpload(input:$input) }",
  "variables": {
    "input": {
      "session_token": "<getToken 拿到的session_token>",
      "upload": null
    }
  }
}
```

_**map**_

```json
{
  "0": [
    "variables.input.upload"
  ]
}
```

_**0**_

```
<上传的文件对象>
```

#### CURL 示例

假设获取到的 Session Token 为 `fb098270f0288507928f572c32c88b40`

```shell
curl localhost:4000/graphql \
  -F operations='{"query":"mutation($input:LocalUploadInput!){ localUpload(input:$input) }","variables": {"input":{"session_token": "fb098270f0288507928f572c32c88b40","upload": null}}}' \
  -F map='{ "0": ["variables.input.upload"] }' \
  -F 0=@a.png
```

## 其他

### 服务器版本信息

Query

```graphql
query{
    status{
        version_name
        version_code
        storage_provider
    }
}
```

Output

```json
{
  "data": {
    "status": {
      "version_name": "0.1.0",
      "version_code": 1,
      "storage_provider": "local"
    }
  }
}
```

Query

```graphql

```

Input

```json

```