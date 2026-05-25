一个基于 Gin + Telegram Bot + OpenIM 的客服桥接系统

实现：

- Telegram 用户与客服实时聊天
- OpenIM 作为客服 IM 后台
- Gin 作为协议桥接层
- 支持 OpenIM Webhook 回调
- 支持多 Telegram 用户会话

---

# 架构

```text
Telegram User
      ↓
Telegram Bot
      ↓
Gin Server
      ↓
OpenIM
      ↓
customer_service