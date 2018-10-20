# Rocket

利用PT站新种免费的规则，在不增加下载流量的前提下刷上传流量

### 支持网站

- [ ] Ourbits --> `ourbits`
- [ ] M-Team --> `mteam`

### 使用方法

运行前请根据自己的需求调整配置文件(`config.json`)中的内容

```json
{
    "transmission": {
        "url": "localhost",
        "user": "admin",
        "password": "12345678"
    },
    "site": {
        "name": "ourbits",
        "passkey": "YOUR_PASSKEY_HERE",
        "cookie": "YOUR_COOKIE_HERE"
    },
    "size_limit": 60000,
    "storage_limit": 500000,
    "policy": "default"
}
```

- transmission: Transmission RPC相关配置
    * url: RPC地址
    * user: RPC用户名
    * password: RPC密码
- site: PT站点相关配置
    * name: 站点名称，参考`支持网站`
    * passkey: 站点passkey，用于识别身份
- size_limit: 种子大小限制，单位`MB`
- storage_limit: 存储空间大小限制，单位`MB`
- policy: 运行策略
    * default: 存储空间已满时，优先删除创建时间较早的种子