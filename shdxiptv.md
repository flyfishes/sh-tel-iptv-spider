# 上海电信 IPTV 外部接口调用顺序

> ⚠️ 本文由 AI 自动生成，仅供参考。

---

## 调用顺序总览

```
① 4K 登录                    GET   /iptv3a/4kLogAuth.do
    ↓ (返回表单 → 提交)
② R4K 登录                   动态  (表单 action URL)
    ↓ (返回表单 + encryptToken → 追加 authenticator 提交)
③ OTT 认证                   动态  (表单 action URL)
    ↓ (返回 channelArray + epgform)
④ 解析频道                    非HTTP (JS 解析 → channels 表)
    ↓
⑤ EPG 首页                   动态  (form#epgform action URL)
    ↓ (获取 UserToken → JS 重定向)
⑥ EPG 负载均衡               GET   (JS 跳转 URL)
    ↓ (返回 EPG 登录表单 → 追加 stbinfo 提交)
⑦ EPG 门户认证 ★             动态  (表单 action URL)
    ↓ (获取 JSESSIONID + EPGHostUrl)
⑧ 门户首页                   GET   {EPGHostUrl}/portal.jsp
    ↓
⑨ Session 检查               GET   {EPGHostUrl}/service/auth/AuthByAjax.jsp?action=auth
    ↓
⑩ 频道列表                   POST  {EPGHostUrl}/function/ajax/epg7getChannelByAjax.jsp
    ↓
⑪ 节目单 ×N                  POST  {EPGHostUrl}/function/ajax/epg7getChannelByAjax.jsp
```

> ⑨ 每 30 分钟检查一次，过期则从 ① 重走认证。

---

## 接口详情

### ① 4K 登录

```
GET http://{auth_host}/iptv3a/4kLogAuth.do
  ?Action=Login
  &UserID={uid}            例: 65059297@etv1
  &SN={sn}                 例: 00030021535101708110051Q
  &Type=iptv4k
  &Mode=MENU.SMG-4K
  &FCCSupport=1

Cookie:    无
UA:        webkit;Resolution(PAL,720P,1080P,2106P,4K)
响应:      HTML（含下一步认证表单 <form>）
```

`auth_host` 默认 `222.68.208.73:7001`，可在 `config.yaml` → `stb.auth_host` 覆盖。

---

### ② R4K 登录

```
URL:       从 ① HTML <form action> 动态拼接
方法:      从 <form method> 动态确定
请求体:    ① 返回表单的全部 <input> 字段
Cookie:    继承 ①
响应:      HTML（含下一步表单 + JS 变量 encryptToken）
```

---

### ③ OTT 认证

```
URL:       从 ② HTML <form action> 动态拼接
请求体:    ② 原有表单字段 + authenticator
Cookie:    继承 ①/②

authenticator = AES-128-ECB(密钥=MD5("123456"), PKCS5填充, JSON原文={
    "Randon":     "{随机}",
    "EncryToken": "{encryptToken from ②}",
    "UserID":     "{uid}",
    "SN":         "{sn}",
    "IP":         "{ip}",
    "MAC":        "{mac}",
    "MagicCode":  "CTC",
    "UpdateTime": "20230301175307"
})

响应:      HTML（含 form#epgform + JS 变量 channelArray）
```

---

### ④ 解析频道

非 HTTP 请求，执行 ③ HTML 中的 JS 解析 `channelArray`，写入本地 `channels` 表。

---

### ⑤ EPG 首页

```
URL:       从 ③ HTML form#epgform action 动态拼接
请求体:    form#epgform 全部表单字段
Cookie:    累积 ①/②/③
提取:      UserToken → 用于 ⑦ 签名
响应:      HTML（含 JS 跳转: top.document.location='...'）
```

---

### ⑥ EPG 负载均衡

```
URL:       从 ⑤ HTML <script> 中提取跳转 URL
方法:      GET
Cookie:    累积
响应:      HTML（含 EPG 门户认证表单）
```

---

### ⑦ EPG 门户认证 ★

```
URL:       从 ⑥ HTML <form action> 动态拼接
方法:      从 <form method> 动态确定
请求体:    ⑥ 原有表单字段 + stbtype + stbinfo

stbtype = {stb.type}           例: B860
stbinfo = RSA私钥签名 PKCS1v15, 无哈希(硬编码1024bit PKCS#8私钥,
           plainData=UserToken第7位插入"37AE")
        → 十六进制大写

Cookie:    累积
响应:      HTML（内嵌JS: jsSetConfig("SessionID","..."), jsSetConfig("IpPort","...")）

提取:      SessionID → JSESSIONID Cookie
           IpPort:Port + framecode → EPGHostUrl
           例: http://218.83.165.40:8084/iptvepg/frame1413/
```

---

### ⑧ 门户首页

```
GET {EPGHostUrl}/portal.jsp
Cookie: JSESSIONID={SessionID}
```

---

### ⑨ Session 状态检查

```
GET {EPGHostUrl}/service/auth/AuthByAjax.jsp?action=auth
Cookie: JSESSIONID={SessionID}

判断: Content-Type 含 "json" → 有效, 否则 → 过期需重走认证
```

---

### ⑩ 频道列表

```
POST {EPGHostUrl}/function/ajax/epg7getChannelByAjax.jsp
Content-Type: application/x-www-form-urlencoded
Cookie: JSESSIONID={SessionID}

action=getChannelList
cateID=000406

响应 JSON:
{
  "result": [{
    "code":    "cctv1hd",
    "name":    "CCTV-1 HD",
    "ID":      "ch00000012",
    "mixNo":   "8",
    "tsTime":  0,
    "isTs":    "1",
    "isCharge":"0"
  }]
}
```

---

### ⑪ 节目单（逐频道，间隔 500ms）

```
POST {EPGHostUrl}/function/ajax/epg7getChannelByAjax.jsp
Content-Type: application/x-www-form-urlencoded
Cookie: JSESSIONID={SessionID}

action=getChannelProg
code={频道code}              例: cctv1hd
channelID={频道ChID}         例: ch00000012
startTime={7天前毫秒时间戳}
endTime={3天后毫秒时间戳}
offset=0
limit=2000

响应 JSON:
{
  "result": [{
    "ID":        "epg_prog_001",
    "name":      "今日说法",
    "startTime":  1751040000000,
    "endTime":    1751041800000,
    "status":    "1",
    "interRecordStatus": "0"
  }]
}
```

> 单频道失败不中断，跳过继续。

---

## 汇总

| # | 端点 | 方法 | 关键参数 |
|---|------|------|----------|
| ① | `/iptv3a/4kLogAuth.do` | GET | Action,UserID,SN,Type,Mode,FCCSupport |
| ② | 动态(form action) | 动态 | 上一步全部表单字段 |
| ③ | 动态(form action) | 动态 | 表单字段 + authenticator (AES加密) |
| ⑤ | 动态(form#epgform action) | 动态 | 表单字段 (含UserToken) |
| ⑥ | 动态(JS跳转) | GET | — |
| ⑦ | 动态(form action) | 动态 | 表单字段 + stbtype + stbinfo (RSA签名) |
| ⑧ | `{EPGHostUrl}/portal.jsp` | GET | — |
| ⑨ | `{EPGHostUrl}/service/auth/AuthByAjax.jsp` | GET | action=auth |
| ⑩ | `{EPGHostUrl}/function/ajax/epg7getChannelByAjax.jsp` | POST | action=getChannelList, cateID=000406 |
| ⑪ | `{EPGHostUrl}/function/ajax/epg7getChannelByAjax.jsp` | POST | action=getChannelProg, code, channelID, startTime, endTime, offset, limit |
