# 数据库表结构

## Users

|     Field     |  Type  | Description |
|:-------------:|:------:|:-----------:|
|      id       | String |     ID      |
|   username    | String |     昵称      |
|    avatar     | String |     头像      |
|   password    | String | 密码(BCrypt)  |
|     email     | String |     邮箱      |
|   subscribe   |  Int   |    订阅类型     |
|  expire_time  |  Int   |    过期时间     |
| register_time |  Int   |    注册时间     |

## Memories

|    Field     |  Type  | Description |
|:------------:|:------:|:-----------:|
|      id      | String |     ID      |
|     uid      | String |     UID     |
|    title     | String |   标题（可空）    |
|   content    |  Json  |     内容      |
|   deleted    |  Bool  |   软删除标志位    |
| created_time | String |    创建时间     |
| update_time  | String |    修改时间     |

## Histories

|    Field     |  Type  | Description |
|:------------:|:------:|:-----------:|
|      id      | String |     ID      |
|     uid      | String |     UID     |
|     ref      | String | 原Memory ID  |
|    title     | String |    历史标题     |
|   content    |  Json  |    历史内容     |
| created_time | String |    创建时间     |

## Comments

|    Field     |  Type  | Description |
|:------------:|:------:|:-----------:|
|      id      | String |     ID      |
|     uid      | String |    用户ID     |
|   content    |  Json  |     内容      |
|   deleted    |  Bool  |   软删除标志位    |
| created_time | String |    创建时间     |

## HashTags

|    Field     |  Type  |             Description             |
|:------------:|:------:|:-----------------------------------:|
|      id      | String |                 ID                  |
|     uid      | String |                用户ID                 |
|     name     |  Json  |                标签名称                 |
|   deleted    |  Bool  | 软删除标志位，软删除后列表不可见，但仍可通过Memories的引用查看 |
| created_time | String |                创建时间                 |

## Resources

|    Field     |  Type  | Description |
|:------------:|:------:|:-----------:|
|      id      | String |     ID      |
|     uid      | String |    用户ID     |
|     path     |  Json  |   资源所在路径    |
|     size     |  Json  | 资源大小（Byte）  |
|   deleted    |  Bool  |   软删除标志位    |
| created_time | String |    创建时间     |

## Subscribes

|    Field     |  Type  | Description |
|:------------:|:------:|:-----------:|
|      id      | String |     ID      |
|     name     | String |    订阅名称     |
|   capacity   |  Int   | 资源额度（Byte）  |
|   deleted    |  Bool  |   软删除标志位    |
| created_time | String |    创建时间     |