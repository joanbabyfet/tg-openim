一个基于 Gin + Telegram Bot + OpenIM 的客服桥接系统

实现：

- Telegram 用户发送消息
- 自动注册 OpenIM 用户
- 转发到 OpenIM 客服
- 客服回复后
- 自动回发 Telegram

---

# 架构流程

```text
Telegram User
      ↓
Telegram Webhook
      ↓
Gin Handler
      ↓
自动注册 OpenIM 用户
      ↓
发送消息到 OpenIM
      ↓
customer_service 收到消息
      ↓
客服回复
      ↓
OpenIM Callback
      ↓
发送回 Telegram
```

---

# 项目结构

```text
project/
├── config/
│   └── config.go
├── handler/
│   ├── openim.go
│   └── tg.go
├── model/
│   ├── mapping.go
│   └── telegram.go
├── service/
│   ├── openim.go
│   ├── telegram.go
│   └── token.go
├── main.go
├── go.mod
└── README.md
```

---

# 用户 ID 设计

Telegram 用户：

```text
tg_<telegram_user_id>
```

例如：

```text
tg_7833973372
```

原因：

- Telegram username 可修改
- Telegram ID 永久不变
- 避免与 App 用户冲突

---

# OpenIM 用户

## Telegram 用户

自动注册：

```text
tg_7833973372
```

---

## 客服账号

固定：

```text
customer_service
```

需提前手动注册。

---

# 注册 customer_service

调用 OpenIM REST API：

```bash
curl -X POST http://127.0.0.1:10002/user/user_register \
-H "Content-Type: application/json" \
-H "token: ADMIN_TOKEN" \
-H "operationID: 123456" \
-d '{
  "users":[
    {
      "userID":"customer_service",
      "nickname":"客服"
    }
  ]
}'
```

---

# 获取 OpenIM Admin Token

默认：

```text
userID: imAdmin
secret: openIM123
```

获取 token：

```bash
curl -X POST http://127.0.0.1:10002/auth/get_admin_token \
-H "Content-Type: application/json" \
-H "operationID: 123456" \
-d '{
  "secret":"openIM123",
  "userID":"imAdmin"
}'
```

---

# Telegram Bot 配置

## 创建 Bot

联系：

```text
@BotFather
```

获取：

```text
BOT_TOKEN
```

---

# 设置 Telegram Webhook

代码：

```go
webhook, _ := tgbotapi.NewWebhook(
    "https://your-domain.com/webhook",
)
```

Telegram 要求：

- HTTPS
- 公网域名
- 443/8443/80/88 端口

不支持：

```text
localhost
127.0.0.1
```

---

# OpenIM Callback 配置

必须配置 OpenIM callback。

例如：

```yaml
callback:
  afterSendSingleMsg:
    enable: true
    callbackURL: "http://host.docker.internal:8080/openim/callback"
```
---

# 启动项目

```bash
go run main.go
```

日志：

```text
Telegram Bot 已启动
服务启动: 8080
```

---

# 测试流程

## 1. Telegram 发消息

```text
hello
```

---

## 2. OpenIM 收到

客服账号：

```text
customer_service
```

收到消息。

---

## 3. 客服回复

```text
你好
```

---

## 4. Telegram 收到

```text
客服回复: 你好
```

---

# 常见问题

---

## 1. OpenIM callback 没触发

原因：

- callback 未配置
- callback URL 错误
- Docker 无法访问宿主机

检查：

```bash
docker logs -f openim-server
```

---

## 2. Telegram 收不到回复

检查：

```go
TgUserMap[msg.RecvID]
```

是否正确。

---

## 3. content.Text 为空

很多 OpenIM 版本：

```json
"content":"{\"content\":\"hello\"}"
```

content 是字符串 JSON。

需要二次解析。

---

## 4. customer_service 无法登录官方 Web

官方 Web 默认：

```text
手机号 + 验证码
```

customer_service 是纯 userID 用户。

建议：

- 自己做客服后台
- 使用 OpenIM SDK 登录

---

# 推荐生产架构

```text
Telegram
    ↓
Gin Gateway
    ↓
OpenIM
    ↓
客服后台
```

建议：

- OpenIM 仅负责 IM
- 用户体系自己维护
- RBAC 自己实现
- OpenIM 只做消息能力

---

# 技术栈

- Golang
- Gin
- Telegram Bot API
- OpenIM
- Docker
- MongoDB

---

# License

[MIT License](https://github.com/joanbabyfet/tg-openim/blob/main/LICENSE)