# API 文档

## 已知问题

Memory 删除时Comments及其子回复未删除（有必要删除）  
Memory 删除时历史记录未删除（有必要删除）  
Memory 更新后存在历史记录对资源的引用(忽略)

## 待测试接口

```
其他
资源链接生成测试
其他用户内容可见性测试
```

以下接口为已经过测试的接口，可查看GraphQL Schema文件夹 [graph/schema](../graph/schema) 查看全部接口

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

### 获取 Memory 详情

Query

```graphql
query ($input: ID!){
    memory(input: $input){
        id
        user{
            id
            username
            avatar
        }
        title
        content
        hashtags{
            id
            name
        }
        archived
        create_time
        update_time
    }
}
```

Input

```json
{
  "input": "644d2af60cef546afdd38365"
}
```

### 获取 Memories

Query

```graphql
query ($input:ListInput!){
    allMemories(input:$input){
        id
        user{
            id
            username
            avatar
        }
        title
        content
        hashtags{
            id
            name
        }
        archived
        create_time
        update_time
    }
}
```

Input

```json
{
  "input": {
    "page": 0,
    "size": 10,
    "byCreate": false,
    "desc": true,
    "archived": false
  }
}
```

### 获取指定话题下的 Memories

Query

```graphql
query ($tag: ID!, $input: ListInput!) {
    allMemoriesByTag(tag: $tag, input: $input) {
        id
        user {
            id
            username
            avatar
        }
        title
        content
        hashtags {
            id
            name
        }
        archived
        create_time
        update_time
    }
}
```

Input

```json
{
  "tag": "6448e4eb03bbc9b5380ded1d",
  "input": {
    "page": 0,
    "size": 10,
    "byCreate": false,
    "desc": true,
    "archived": false
  }
}
```

### 更新 Memory

Query

```graphql
mutation ($input:UpdateMemoryInput!){
    updateMemory(input:$input)
}
```

Input
`title`与`content`非必须

```json
{
  "input": {
    "id": "6448e3654c750ecc11b96999",
    "title": "<新标题>",
    "content": "<新内容> 附加资源采用${资源ID}引用，获取Memory时会自动替换为资源链接![Hello](${644d3f61dd12d1c7d2d65b5c})"
  }
}
```

### 归档 Memory

Query

```graphql
mutation ($input:ID!,$archived:Boolean! = true){
    archiveMemory(input:$input,archived:$archived)
}
```

Input

```json
{
  "input": "6448e3654c750ecc11b96999",
  "archived": true
}
```

### 删除 Memory

Query

```graphql
mutation ($input:ID!){
    deleteMemory(input:$input)
}
```

Input

```json
{
  "input": "6448e3654c750ecc11b96999"
}
```

## 历史记录

历史记录仅保留最基础的文本信息，资源删除后无法通过历史记录恢复  
Query

```graphql
query ($id: ID!, $page: Int64!, $size: Int64!, $desc: Boolean!) {
    allHistories(
        id: $id
        page: $page
        size: $size
        desc: $desc
    ) {
        id
        user {
            id
            username
            avatar
        }
        title
        content
        hashtags {
            id
            name
        }
        create_time
    }
}
```

Input

```json
{
  "id": "644d2af60cef546afdd38365",
  "page": 0,
  "size": 10,
  "desc": true
}
```

## 标签

### 标签列表

Query

```graphql
query($input:ListInput!){
    allHashTags(input:$input){
        id
        name
        archived
        create_time
        update_time
    }
}
```

Input

```json
{
  "input": {
    "page": 0,
    "size": 10,
    "byCreate": false,
    "desc": true,
    "archived": false
  }
}
```

### 更新标签

Query

```graphql
mutation($input:HashTagInput!){
    updateHashTag(input:$input)
}
```

Input  
`name`和`archived`字段非必须

```json
{
  "input": {
    "id": "644d2af60cef546afdd38364",
    "name": "<新标签名>",
    "archived": true
  }
}
```

### 删除标签

**仅已归档标签且无Memory引用的标签可被删除**  
Query

```graphql
mutation($input:ID!){
    deleteHashTag(input:$input)
}
```

Input

```json
{
  "input": "644d2af60cef546afdd38363"
}
```

## 回复

### 新增回复

Query

```graphql
mutation($input:AddCommentInput!){
    addComment(input:$input)
}
```

Input

```json
{
  "input": {
    "id": "<subComment为false时为回复的Memory id,true时为回复对象的id>",
    "subComment": false,
    "content": "<回复内容>"
  }
}
```

### 更新回复

Query

```graphql
mutation($input:UpdateCommentInput!){
    updateComment(input:$input)
}
```

Input  
`content`和`archived`字段非必须

```json
{
  "input": {
    "id": "644d318b8aa787230d087d24",
    "content": "<修改后的回复>",
    "archived": false
  }
}
```

### 删除回复

**仅已归档回复可被删除**  
Query

```graphql
mutation($input:ID!){
    deleteComment(input:$input)
}
```

Input

```json
{
  "input": "644d30cf5a41ee83413e3301"
}
```

### 回复列表

Query

```graphql
query ($id: ID!, $page: Int64!, $size: Int64!, $desc: Boolean!) {
    allComments(
        id: $id
        page: $page
        size: $size
        desc: $desc
    ) {
        id
        user {
            id
            username
            avatar
        }
        memory{
            id
            content
        }
        content
        hashtags {
            id
            name
        }
        create_time
    }
}
```

Input

```json
{
  "id": "<Memory ID>",
  "page": 0,
  "size": 10,
  "desc": true
}
```

### 子回复列表

Query

```graphql
query ($id: ID!, $page: Int64!, $size: Int64!, $desc: Boolean!) {
    subComments(
        id: $id
        page: $page
        size: $size
        desc: $desc
    ) {
        id
        user {
            id
            username
            avatar
        }
        comment{
            id
            content
        }
        content
        hashtags {
            id
            name
        }
        create_time
    }
}
```

Input

```json
{
  "id": "<Top-Level 回复 ID>",
  "page": 0,
  "size": 10,
  "desc": true
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

## 资源

### 获取资源列表

Query

```graphql
query ($page: Int64!, $size: Int64!, $byCreate: Boolean!, $desc: Boolean!) {
    allResources(page: $page, size: $size, byCreate: $byCreate, desc: $desc) {
        id
        user {
            id
            username
            avatar
        }
        path
        size
        memories{
            id
            title
        }
        create_time
    }
}
```

Input

```json
{
  "page": 0,
  "size": 10,
  "byCreate": false,
  "desc": true
}
```

### 删除资源

只有没有Memory引用改资源时才可删除  
Query

```graphql
mutation ($input:ID!){
    deleteResource(input:$input)
}
```

Input

```json
{
  "input": "644d43d6b3f18106402b5c6c"
}
```

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

## 订阅相关

### 查看所有可用订阅

Query

```graphql
{
    allSubscribes {
        id
        name
        capacity
        available
        create_time
        update_time
    }
}
```

### 创建订阅

Query

```graphql
mutation($input:AddSubscribeInput!){
    addSubscribe(input:$input)
}
```

Input

```json
{
  "input": {
    "name": "1GB",
    "capacity": 1024000000,
    "enable": true
  }
}
```

### 更新订阅

Query

```graphql
mutation($input:UpdateSubscribeInput!){
    updateSubscribe(input:$input)
}
```

Input

```json
{
  "input": {
    "id": "644dddce8a942a83b12f96ca",
    "name": "1GB",
    "capacity": 1024000000,
    "enable": true
  }
}
```

### 删除订阅

当且仅当该订阅下的用户数为0时允许删除

Query

```graphql
mutation($input:ID!){
    deleteSubscribe(input:$input)
}
```

Input

```json
{
  "input": "644dddce8a942a83b12f96ca"
}
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