# 配置说明

::: tip 提示
本文档旨在针对配置文件中的各个字段进行详细介绍
:::

## 数据库配置

::: tip 提示  
使用默认提供的 docker-compose.yml 部署时无需对其进行修改
:::

```yaml
db: # 数据库配置
  mongo_addr: "mongo:27017" # MongoDB 数据库地址
  mongo_db: "memox" # MongoDB 数据库名称
  redis_addr: "redis:6379" # Redis 地址
  redis_db: 0 # Redis 数据库编号
```

## 调试模式

```yaml
debug: true # 调试模式
```

**调试模式开启后**

- Gin会在日志中输出每一次请求
- 会开启 Playground

::: tip 如果想要单独控制 Gin 的日志模式，请设置环境变量
最少日志，不输出请求信息

```text
GIN_MODE=release
```

较少日志，输出请求信息

```text
GIN_MODE=test
```

最多日志，输出 Gin 调试信息和请求信息 (不推荐)

```text
GIN_MODE=debug
```

:::

## 用户配置

```yaml
user: # 用户配置
  token_expire: 43200 # 30 * 24 * 60 = 30 天
  token_secret: "<JWT密钥>" # 不要使用默认密钥
```

::: tip token_expire  
登录有效时间(分钟)，超过设定的时间后会要求重新登录
:::

::: tip JWT密钥
JWT 密钥对格式没有特殊要求，可使用密码生成器生成高强度密钥，如

```text
76vL#o0RtZVgyr*ZUasq5qB@28y*byFu
```

但请不要在生产环境使用以上示例密钥
:::

## 邮件发送配置

```yaml
mail: # 邮件发送配置
  code_expire: 10 # 验证码有效时长(分钟)
  code_length: 6 # 验证码长度
  code_cool_down: 1 # 验证码冷却时间(分钟)
  smtp_mail_host: "smtp.exmail.qq.com" # SMTP服务器地址
  smtp_mail_port: 465 # SMTP服务器端口
  smtp_mail_user: "example@example.com" # 邮箱用户名
  smtp_mail_pwd: "<邮箱密码>" # 邮箱密码
  smtp_mail_nickname: "忆盒" # 发件人昵称
  subject: "【忆盒】邮箱验证提醒" # 邮件主题
  template: "/app/env/template.html" # 邮件模板
```

::: danger 邮件模板
邮件模板的路径在使用 `docker-compose` 进行部署时请不要修改，实际路径应在 `docker-compose.yml`
所在目录下的 `env/template.html`  
具体请参考 [docker-compose.yml](https://github.com/MemoxApp/server/blob/main/docker-compose.yml) 文件中的 `volumes`
数据卷映射关系
:::

## 订阅配置

配置新用户注册后默认分配的订阅计划

```yaml
subscribe: # 订阅配置
  default_capacity: 104857600 # 默认容量(Bytes) = 100 MB
  default_subscribe_name: "免费版" # 默认订阅名称
```

## 存储配置

### 本地存储

将资源文件存储在本地，由Server提供静态资源服务，修改配置文件中 `storage` 项：

```yaml
storage: # 存储配置
  storage_provider: "local" # 本地提供资源存储服务
  local: # 本地存储配置
    folder: "/app/data/storage" # 存储文件夹，使用 docker-compose 时不要修改该项
    host: "localhost:8080" # 本地存储访问地址
    schema: "http" # 本地存储访问协议
```

### 百度云对象存储

存储桶所在地区需要支持[事件通知](https://cloud.baidu.com/doc/BOS/s/kjwvyr7st),目前仅支持北京、苏州、广州地区

#### 仅使用百度云对象存储

```yaml
storage: # 存储配置
  storage_provider: "bce" # 百度云
  bce: # 百度云配置
    access_key_id: "<AccessKeyID>" # 百度云AccessKeyID
    secret_access_key: "<SecretAccessKey>" # 百度云SecretAccessKey
    # 如果通过CDN接入此处填写CDN域名(需要包含协议头)
    end_point: "https://su.bcebos.com" # 存储桶Endpoint: <region>.bcebos.com，如需使用 https 协议请添加 https 协议头，如：https://<region>.bcebos.com
    bucket_name: "memox" # 此处填写存储桶名称，需要手动预先创建
    cdn: false # 不使用 CDN
    region: "su" # 存储桶所在区域
    callback_token: "<回调鉴权密钥>" # 存储桶回调通知鉴权密钥，防止伪造请求，尽可能只包含数字与字母避免影响URL参数的解析
```

##### 创建存储桶事件通知

| 名称     | 配置                                                |
|--------|---------------------------------------------------|
| 状态     | **开启**                                            |  
| 名称     | **MemoxUploadNotify**（随意填写）                       |
| 规则ID   | **UploadNotify** （随意填写）                           |
| 产品ID   | **Memox**（随意填写）                                   |
| 加密鉴权   | **关闭**                                            |
| 事件监听配置 | **开启**                                            |
| 监测事件   | 勾选 **PutObject**                                  |
| 覆盖资源   | **前后缀**                                           |
| 资源前缀   | **resources/**                                    |
| 资源后缀   | 留空                                                |
| 触发应用   | **自定义应用**                                         |
| 应用地址   | `http(s)://<host:port>/notify/bce?token=<回调鉴权密钥>` |

::: tip 提示：应用地址
可使用[在线测试HTTP回调服务](https://hooks.upyun.com)多添加一个应用地址用于测试事件触发是否可用
:::

#### 百度云CDN配合对象存储使用

CDN 创建时源站选择 BOS Bucket，主源站 Bucket 选择预先创建的存储桶，域名级回源HOST选择加速域名

```yaml
storage: # 存储配置
  storage_provider: "bce" # 百度云
  bce: # 百度云配置
    access_key_id: "<AccessKeyID>" # 百度云AccessKeyID
    secret_access_key: "<SecretAccessKey>" # 百度云SecretAccessKey
    end_point: "https://cdn.ts.example.com" # CDN 加速地址，必须包含协议头
    cdn: true # 是否通过CDN接入
    cdn_auth_key: "<鉴权密钥，与CDN控制台处填写的一致>" # CDN鉴权密钥
    region: "su" # 存储桶所在区域
    callback_token: "<回调鉴权密钥>" # 存储桶回调通知鉴权密钥，防止伪造请求
```

##### 创建存储桶事件通知

[参见上一步的创建存储桶事件通知](#创建存储桶事件通知)

##### 配置 CDN 鉴权

CDN 配置 > 访问控制 > 高级鉴权  
|名称|配置|
|-----|-----|
|鉴权配置| **开启**|  
|类型选择|**类型B**|  
|主KEY| **配置文件中配置的鉴权密钥cdn_auth_key**|  
|备KEY|留空|  
|时间格式| **十进制**|  
|有效时间| **900** 秒 (根据个人偏好设置)|

##### 配置 CDN 回源权限

CDN 配置 > 回源配置 > 私有Bucket回源 开启  
之后在 对象存储 Bucket 的配置管理中设置 `Bucket权限配置` 状态应该为 `自定义`，展开后点击右边的编辑后编辑自定义权限
|名称|配置|
|-----|-----|
|用户授权|**所有用户**|  
|授权效果|**允许**|  
|权限设置|**全部勾选**|  
|权限高级设置|**关**|  
|资源|**包含**|  
|自定义资源|**关**|  
|条件|**不勾选**|

点击**修改权限**保存