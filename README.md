<center>
<h1>Time Speak Opensource Server</h1>

悄语时光，私语心声。  
Talk to time, talk to self.

<strong>项目当前处于前期开发阶段，敬请期待</strong>
</center>

# 简介

在时光语中，您可以记录自己的想法、经历和感受，您可以随时随地回顾过去的自己，并在当下做出对过去的回复与评价。

时光语中除非您主动分享，否则只有您自己可以访问和查看您所有记录的内容。您可以自由地在App中发表和编辑帖子，也可以随时删除或归档帖子。时光语将帮助您更好地记录和管理自己的想法和回忆。

时光语，意喻时间和回忆，代表与过去的自己对话。时光流淌，一切皆有答案。

# 部署

## Docker Compose

克隆本仓库或复制仓库中的 `docker-compose.yml` 和 `env` 文件夹，根据情况修改[env/example.yaml](env/example.yaml)的配置
完成后执行运行服务

```shell
docker-compose up -d
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
    bucket_name: "timespeak" # 此处填写存储桶名称，需要手动预先创建
    cdn: false # 不使用 CDN
    region: "su" # 存储桶所在区域
    callback_token: "<回调鉴权密钥>" # 存储桶回调通知鉴权密钥，防止伪造请求，尽可能只包含数字与字母避免影响URL参数的解析
```

##### 创建存储桶事件通知

状态： **开启**
名称： **TimeSpeakUploadNotify**（随意填写）
规则ID： **UploadNotify** （随意填写）
产品ID： **TimeSpeak**（随意填写）
加密鉴权：**关闭**
事件监听配置： **开启**
监测事件：勾选`PutObject`
覆盖资源： **前后缀** 资源前缀:`resources/`,资源后缀留空
触发应用： **自定义应用**
应用地址：http://<host:port>/notify/bce?token=<回调鉴权密钥>
同时可使用[在线测试HTTP回调服务](https://hooks.upyun.com)多添加一个应用地址用于测试事件触发是否可用

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

同上

##### 配置 CDN 鉴权

CDN 配置 > 访问控制 > 高级鉴权
鉴权配置： **开启**
类型选择：**类型B**
主KEY： 配置文件中配置的鉴权密钥cdn_auth_key
备KEY：留空
时间格式： **十进制**
有效时间： 900 秒 (根据个人偏好设置)

##### 配置 CDN 回源权限

CDN 配置 > 回源配置 > 私有Bucket回源 开启
之后在 对象存储 Bucket 的配置管理中设置 `Bucket权限配置` 状态应该为 `自定义`，展开后点击右边的编辑后编辑自定义权限

用户授权：**所有用户**
授权效果：**允许**
权限设置：**全部勾选**
权限高级设置：关
资源：**包含**
自定义资源：**关**
条件：**不勾选**

点击修改权限保存

# API 及开发文档

详见 [docs/API.md](docs/API.md)

# LICENSE

```
MIT License

Copyright (c) 2023 Time Speak App

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```