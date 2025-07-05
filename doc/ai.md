你是一名资深全栈工程师，请为我生成一个「轻量级 Nginx 日志查看器」，技术栈 Golang + Vue3 + Element-Plus；日志格式使用下方自定义的 **log_format security**。要求如下：

─────────────────────────────────────────
A. 日志格式（务必按此解析）
示例行（换行仅为示意；实际同一行）：
192.0.2.10:51324 - example.com:443 [10/Jun/2024:12:34:56 +0800] \
"GET /index.html HTTP/1.1" 200 1234 456 rt=0.123 \
uct=0.001 urt=0.120 ust=200 ua=192.0.2.20:9000 \
"https://ref.example.com/" "Mozilla/5.0 ..." TLSv1.3 TLS_AES_256_GCM_SHA384 \
rid=3c6e0b8a9c15224a

字段含义与顺序：
1) client_ip          → 192.0.2.10
2) client_port        → 51324
3) x_forwarded_for    → -                           （可能为空 / "-"）
4) host               → example.com
5) server_port        → 443
6) time_local         → 10/Jun/2024:12:34:56 +0800
7) request            → "GET /index.html HTTP/1.1"  （需拆 Method / Path / HTTPVer）
8) status             → 200
9) body_bytes_sent    → 1234
10) request_length    → 456
11) request_time      → rt=0.123
12) upstream_connect  → uct=0.001                  （空时为 "-")
13) upstream_time     → urt=0.120                  （空时为 "-")
14) upstream_status   → ust=200                    （空时为 "-")
15) upstream_addr     → ua=192.0.2.20:9000         （空时为 "-")
16) http_referer      → "https://ref.example.com/"
17) http_user_agent   → "Mozilla/5.0 ..."
18) ssl_protocol      → TLSv1.3                    （若 HTTP 则为 "-"）
19) ssl_cipher        → TLS_AES_256_GCM_SHA384     （若 HTTP 则为 "-"）
20) request_id        → rid=3c6e0b8a9c15224a       （1.17+ 可用，缺省时为 "-"）

请在 Go 端写出解析正则或 strings.Split 逻辑，最终映射到如下结构体：

type LogRow struct {
Time         time.Time
ClientIP     string
ClientPort   string
XFF          string
Host         string
ServerPort   string
Method       string
Path         string
HTTPVer      string
Status       int
BodyBytes    int64
ReqBytes     int64
ReqTime      float64
UpConnTime   sql.NullFloat64   // 允许 NULL
UpRespTime   sql.NullFloat64
UpStatus     sql.NullInt64
UpAddr       sql.NullString
Referer      string
UA           string
TLSProto     string
TLSCipher    string
ReqID        string
}

B. 整体功能
1. 目录 /var/log/nginx
    - 当前运行中的 access.log（未压缩）
    - 最近 N(默认31) 个 access.log.*.gz
2. 每次调用 /api/logs 时：
   a) 枚举文件列表（access.log + 最近 N 天 .gz）
   b) 若文件 mtime 与内存缓存一致 → 直接用缓存
   否则 → 流式解压 & 解析 → 保存缓存
   c) 在内存切片中过滤 / 分页 / 排序后返回 JSON
3. 过滤参数：
   ip, status, path, from(起始时间), to(结束时间), limit, offset
4. BasicAuth ；账号密码写在 config.yaml
5. 前端：Vue3 + Element-Plus 表格
    - 过滤栏：IP / 状态码 / Path / 时间范围
    - 彩色 4xx/5xx
    - 分页 200/页
    - 导出 CSV（当前过滤结果）
6. 可选 Stats 页面：近 31 天 PV/UV、Top IP（简单 JS 侧统计即可）

C. 目录结构
nginx-viewer/
├── server/
│   ├── main.go
│   ├── config.go
│   ├── parser.go      // 解析 security 格式
│   ├── cache.go       // 文件 mtime + []LogRow
│   ├── api.go
│   └── util.go
└── web/
├── vite.config.ts
└── src/
├── main.ts
├── App.vue
├── views/Logs.vue
└── components/FilterBar.vue

D. 输出顺序
1) 目录树说明
2) server 关键文件代码（main.go→parser.go→cache.go→api.go）
3) web 关键文件代码（main.ts→Logs.vue 等）
4) sample config.yaml
5) README：编译 & 运行指令；logrotate 示例
   ─────────────────────────────────────────

注意事项：
• parser.go 中给出完整正则 (或 strings.Index/Trim)；确保能正确处理可能为 “-” 的字段。  
• 对 .gz 文件使用 gzip.NewReader；对于当前 access.log 直接用 os.Open。  
• 缓存 key = filepath + mtime；解析结果存 [][]LogRow，避免重复解压。  
• 在过滤 / 排序之后才做分页，保证准确结果。

请严格按上述条目输出项目代码与说明，并保证代码能直接 go build && pnpm build 通过。



