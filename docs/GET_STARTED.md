# 快速开始

克隆[本仓库](https://github.com/MemoxApp/server)

```shell
git clone https://github.com/MemoxApp/server.git
```

::: tip 提示
您也可以仅下载仓库中的 `docker-compose.yml` 文件和 `env` 文件夹，但需要保证文件组织结构与仓库相同
:::

参考 [配置说明](CONFIG.md)
根据需要修改[env/example.yaml](https://github.com/MemoxApp/server/blob/main/env/example.yaml)的配置
完成后执行运行服务

```shell
docker-compose up -d
```

::: tip 权限说明
第一个注册的用户拥有管理员权限，其余用户均为普通用户
:::