<div align="center">

<img src="https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go" alt="Go">
<img src="https://img.shields.io/badge/License-MIT-green?style=flat" alt="License">
<img src="https://img.shields.io/badge/Platform-Linux%20|%20Windows%20|%20OpenWrt-blue?style=flat" alt="Platform">

</div>

# 📺 sh-tel-iptv-spider

![](demo.gif)

上海电信 IPTV 抓取程序 —— 自动抓取 **EPG 节目单** 与 **M3U8 播放地址**，并写入 MySQL。

---

## ⚠️ 使用须知

> 🔒 **请务必遵守以下规则，否则将停止维护。**
   本软件完全免费，如果你付费购买，请告诉我。并且立即退款
- ❌ **禁止商业化**：不得用于闲鱼、公司等盈利服务，包括代安装、代理服务
- ❌ **禁止宣传**：不得在任何平台（小红书、论坛、QQ 群等）宣传本项目或贴链接
- 🤫 **低调使用，请勿张扬**
- 📌 唯一授权发布：**恩山论坛 - 公子薛**
- 📌 使用本程序需要一定技术水平，伸手党、白痴问题一律无视

---

## 📋 环境要求

| 依赖 | 说明 |
|------|------|
| 上海电信 IPTV 机顶盒 | 已开通 IPTV 服务，获取机顶盒账号 |
| MySQL 数据库 | 存储频道、EPG、认证信息 |
| IPTV 专网访问 | 需解决路由，确保能访问专网 |

> ⚠️ **重要**：程序必须在能访问 IPTV 专网的环境运行，公网无法抓取。回放地址与权健绑定，仅限本人使用。

---

## 🚀 快速开始

### 下载二进制

| 文件 | 平台 |
|------|------|
| `sh-tel-iptv-spider_linux_386` | Linux x86 32位 |
| `sh-tel-iptv-spider_linux_amd64` | Linux x86 64位 |
| `sh-tel-iptv-spider_linux_arm` | Linux ARM 32位 |
| `sh-tel-iptv-spider_linux_arm64` | Linux ARM 64位 |
| `sh-tel-iptv-spider_windows_386.exe` | Windows 32位 |
| `sh-tel-iptv-spider_windows_amd64.exe` | Windows 64位 |

### OpenWrt 用户

```bash
# 安装时区数据（OpenWrt 默认缺少）
opkg update
opkg install zoneinfo-asia
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 运行
./iptv-spider-sh
```

### OpenWrt 后台常驻（procd）

创建服务脚本 `/etc/init.d/iptv-spider`：

```bash
#!/bin/sh /etc/rc.common

START=99
USE_PROCD=1

PROG=/root/sh-tel-iptv-spider_linux_amd64
CONFIG=/root/config.yaml

start_service() {
    procd_open_instance
    procd_set_param command "$PROG" -c "$CONFIG"
    procd_set_param respawn          # 进程挂了自动重启
    procd_set_param stdout 1         # 输出到 logread
    procd_set_param stderr 1
    procd_close_instance
}
```

然后：

```bash
chmod +x /etc/init.d/iptv-spider
/etc/init.d/iptv-spider enable    # 开机自启
/etc/init.d/iptv-spider start     # 立即启动
```

常用管理命令：

```bash
/etc/init.d/iptv-spider status    # 查看运行状态
/etc/init.d/iptv-spider restart   # 重启服务
/etc/init.d/iptv-spider stop      # 停止服务
logread | grep iptv-spider        # 查看日志
```

### 配置文件

编辑 `config.yaml`，填入 MySQL 连接信息、IPTV 认证参数以及自定义频道映射，启动即可。

---

## ⚙️ 功能说明

- **语言**：Go，编译为单一可执行文件
- **跨平台**：Linux / OpenWrt / Windows

| 功能 | 说明 |
|------|------|
| 📡 频道列表 | 自动抓取 IPTV 频道信息 |
| 📅 EPG 节目单 | 抓取节目数据并入库 |
| 🎬 M3U8 地址 | 生成播放列表，支持自定义频道映射 |
| 🗄️ 数据持久化 | 全部写入 MySQL 数据库 |
| 🌐 Web 监控 | 内置状态页面，支持在线查看频道、下载 M3U8/EPG |
| 📡 API 接口 | 完整 REST API，详见 [API.md](API.md) |

---

## 📂 数据库结构

| 表名 | 说明 |
|------|------|
| `auth_infos` | 认证权健存储 |
| `channel_infos` | 频道列表 |
| `channels` | 频道源地址 |
| `epg_details` | 节目单详情 |
| `m3u8_mappings` | 频道分组映射 |

> 建表 SQL 请查看源码。

---

## ⚠️ 限制

- ✅ 仅支持 **上海电信 IPTV**
- ❌ 不支持其他地区电信 / 联通 / 移动

---

## 🤝 贡献

代码写得比较随意 😅，欢迎：

- Fork 项目
- 提交 PR
- 提出 Issue

---

## 📄 免责声明

1. 本程序仅供学习与研究使用
2. 禁止用于商业用途
3. 使用本程序产生的任何法律问题与作者无关
4. 使用即表示同意自行承担风险
