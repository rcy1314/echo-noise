# MCP 接入（AI 客户端）

### 连接方式与配置建议

- 重要：使用 `docker exec` 方式需要本地已安装并运行 Docker，且当前 CLI 指向目标 Docker 守护进程（本机或远程）。如需操作云服务器上的容器，必须先切换到远程 Docker 上下文（`docker context use <remote>`）或配置 `DOCKER_HOST`，否则命令会作用于本机。
- 容器需为 `final-mcp` 目标并已运行，容器名与配置保持一致（例如 `Ech0-Noise`）。
- 环境变量值必须是纯字符串，不要包含反引号或多余空格（如 `NOTE_HOST=http://<服务器IP>:1314`）。
- 不使用本地 Docker 时，可改用 SSH + Stdio 或 HTTP/SSE 网关；详见下文“本地连接远程 MCP（多种方式）”。
- 避免二次占用 `1315` 端口导致 `EADDRINUSE`
- 使用 `docker exec` 启动 `stdio` 握手的实例时，将 `NOTE_HTTP_PORT` 设为 `0`，只进行握手，不再监听 HTTP

```json
{
  "mcpServers": {
    "ech0-noise-mcp-stdio": {
      "command": "docker",
      "args": [
        "exec",
        "-i",
        "-e", "NOTE_HOST=http://localhost:1314",
        "-e", "NOTE_HTTP_PORT=0",
        "-e", "NOTE_TOKEN=你的后台TOKEN",
        "Ech0-Noise",
        "node",
        "/app/mcp/server.bundle.mjs"
      ]
    }
  }
}
```

- 保留 SSE 监控仅用于事件订阅，不作为握手连接

```json
{
  "mcpServers": {
    "ech0-noise-mcp-sse-monitor": {
      "type": "sse",
      "url": "http://localhost:1315/mcp/sse"
    }
  }
}
```

### SSE 握手与保活
- 连接到 `GET /mcp/sse` 后，服务会立即推送握手事件：
  - `event: mcp_hello` 携带服务名称与版本
  - `event: mcp_tools` 携带可用工具列表
  - `event: keepalive` 每 30 秒推送一次，保持连接活跃
- 示例：
  ```bash
  curl -N http://localhost:1315/mcp/sse | head -n 10
  # 预期输出包含：
  # event: mcp_hello
  # data: {"name":"ech0-noise-mcp","version":"0.1.0"}
  # event: mcp_tools
  # data: ["search","publish",...]
  ```
- 提示：SSE 侧仅提供握手信号与运行事件；完整的 MCP 协议交互仍通过 `stdio` 完成。

### 工具与命令

- 工具名称（英文优先）：`search`、`publish`、`delete`、`update`、`message`、`page`、`status`、`calendar`、`config`、`login`、`token`、`rss`
- 中文名称兼容但可能触发客户端校验警告：`搜索`、`发布`、`删除`、`更新`、`消息`、`页面`、`状态`、`日历`、`配置`、`登录`、`令牌`、`RSS`

#### 认证要求与提示
- 无需认证：`search`、`page`、`message`、`status`、`calendar`、`config`、`rss`
- 支持令牌或会话：`publish`
- 令牌或会话（需后端启用 token 路由）：`update`、`delete`、`pin`、`settings`
- 未登录或无令牌时工具会直接返回友好提示：
  - `需要登录或令牌：请先调用 登录 工具，或在配置中设置 NOTE_TOKEN。发布支持令牌；更新/删除/置顶/设置支持令牌（需后端启用）或会话认证。`
  - 令牌无效：`令牌无效或已过期：请在后台重新生成 token，或先 登录。`
  - 后端未启用：`后端未启用 token 路由：请更新后端并重启服务，或使用 登录 获取会话后再操作。`

#### 快捷别名（自然语言更友好）
- 发布别名：`笔记`、`说说`、`说说笔记`（均等价于 `发布`）
- 示例：
  ```json
  { "tool": "笔记", "params": { "content": "这是一段内容", "private": false } }
  { "tool": "说说", "params": { "type": "image", "image": "https://example.com/a.jpg", "content": "配文可选" } }
  ```

#### 避免未格式化的输出
- 一些客户端会原样展示工具调用日志（如 `<tool_use_result>` 等）。建议在提示语中明确要求：
  - “只输出解析后的中文列表，不展示工具调用日志或原始 JSON，不使用任何未渲染标签。”
  - “按中文列出 id、用户名、时间与内容摘要。”

#### 搜索空参回退与输出格式
- 当 `搜索` 未提供 `keyword/query` 时自动回退到分页列表（第一页），避免校验错误
- `搜索/页面` 会同时返回：
  - 可读摘要文本（`id/用户名/时间/内容前200字`）
  - 原始 JSON（便于程序消费）

#### HTTP 调试端点
- 列出工具：
  ```bash
  curl http://localhost:1315/mcp/tools
  ```
- 流式调用工具：返回换行分隔 JSON 事件（`start`、`data`、`end`、`error`）
  ```bash
  curl -N -X POST http://localhost:1315/mcp/tool/search \
    -H 'Content-Type: application/json' \
    -d '{"keyword":"#CDN","page":1,"pageSize":10}'
  ```
- 订阅事件（工具执行的 `tool_start`/`tool_end` 会推送到 SSE）：
  ```bash
  curl -N http://localhost:1315/mcp/sse
  ```

#### 调用示例（AI 客户端/HTTP 皆适用的参数结构）
- 搜索（支持 `keyword` 或 `query`）：
  ```json
  { "tool": "search", "params": { "query": "welcome", "page": 1, "pageSize": 10 } }
  ```
- 分页列表（参数别名兼容）：
  ```json
  { "tool": "page", "params": { "page": 1, "pageSize": 10 } }
  ```
- 获取消息：
  ```json
  { "tool": "message", "params": { "id": "123" } }
  ```
- 发布（纯文本）：
  ```json
  { "tool": "publish", "params": { "type": "text", "content": "一段文本内容" } }
  ```
- 发布（Markdown）：
  ```json
  { "tool": "publish", "params": { "type": "markdown", "content": "# 标题\n正文..." } }
  ```
- 发布（单图）：
  ```json
  { "tool": "publish", "params": { "type": "image", "image": "https://example.com/image.jpg" } }
  ```
- 发布（多图混合，附可选配文）：
  ```json
  { "tool": "publish", "params": { "type": "multipart", "images": ["https://a.jpg","https://b.jpg"], "content": "配文可选" } }
  ```
- 发布（兼容旧字段）：
  ```json
  { "tool": "publish", "params": { "type": "image", "imageURL": "https://example.com/image.jpg" } }
  ```

#### HTTP 调试与事件
- 列出工具：
  ```bash
  curl http://localhost:1315/mcp/tools
  ```
- 流式调用工具（返回换行分隔 JSON 事件）：
  ```bash
  curl -N -X POST http://localhost:1315/mcp/tool/search \
    -H 'Content-Type: application/json' \
    -d '{"query":"welcome","page":1,"pageSize":10}'
  ```
- 订阅事件与握手：
  ```bash
  curl -N http://localhost:1315/mcp/sse
  ```
- 删除：
  ```json
  { "tool": "delete", "params": { "id": "123" } }
  ```
- 更新：
  ```json
  { "tool": "update", "params": { "id": "123", "content": "更新后的文本" } }
  ```
- 置顶消息：
  ```json
  { "tool": "pin", "params": { "id": "123", "pinned": true } }
  ```
  取消置顶：
  ```json
  { "tool": "pin", "params": { "id": "123", "pinned": false } }
  ```
- 系统设置：
  ```json
  { "tool": "settings", "params": { "allowRegistration": true } }
  ```
- 状态与日历：
  ```json
  { "tool": "status", "params": {} }
  { "tool": "calendar", "params": {} }
  ```
- 配置：
  ```json
  { "tool": "config", "params": {} }
  ```
- 登录（会话认证，成功后内部保存 `set-cookie`）：
  ```json
  { "tool": "login", "params": { "username": "admin", "password": "admin" } }
  ```
- 令牌（使用会话或 Token 调用）：
  ```json
  { "tool": "token", "params": {} }
  ```
- RSS：
  ```json
  { "tool": "rss", "params": {} }
  ```

#### 认证与环境变量
- Token 认证：在 MCP 进程环境中设置 `NOTE_TOKEN`，或在 HTTP 调试请求头中附带 `Authorization`
- 会话认证：先调用 `login` 获取服务端 `set-cookie`，内部自动携带会话调用需要认证的工具
- 后端地址：`NOTE_HOST`，容器默认 `http://localhost:1314`（同域接口位于 `/api`）

### 常见问题与排查
- `Invalid content type, expected text/event-stream`：反向代理未正确转发 `/mcp/sse` 或开启缓冲。确保保留 `Content-Type: text/event-stream`，关闭 `proxy_buffering`。
- `EADDRINUSE :::1315`：容器内已有 MCP 监听 `1315`。使用 `docker exec` 启动第二实例时设置 `NOTE_HTTP_PORT=0`，仅进行 `stdio` 握手。
- 工具名校验警告（中文名）：优先使用英文名称 `search/publish/delete/update/message/page/pin/settings/status/calendar/config/login/token/rss`。
- 入参校验错误：
  - `search` 支持 `keyword` 或 `query`
  - `page` 支持 `page/page_number` 与 `pageSize/page_size`
  - `publish` 支持 `type: text|markdown|image|multipart`，图片参数兼容 `image/images/imageURL`
  - 客户端显示“乱码”或原始 JSON：工具默认返回 JSON 文本。请在提示词中要求模型“解析工具返回的 JSON，并用中文列出 id、用户名、时间与内容摘要”，避免直接原样输出。

### 功能概述
- 支持工具：搜索/发布/删除/更新/消息/页面/置顶/设置/状态/日历/配置/登录/令牌/RSS。
- 认证：API Token 或 用户名/密码（会话 Cookie）。
- 传输：标准 Stdio（大多数 MCP 客户端通用）。

### 本地运行
1. 进入 `mcp` 目录并安装依赖：
   ```bash
   cd mcp
   npm install
   ```
2. 使用 Token 认证启动：
   ```bash
   NOTE_HOST=https://note.noisework.cn \
   NOTE_TOKEN=你的_token \
   npm start
   ```
3. 使用用户名密码登录（会话认证）：
   ```bash
   NOTE_HOST=https://note.noisework.cn npm start
   # 在客户端调用“登录”工具：
   # { "username": "admin", "password": "your_password" }
   ```

### Docker 运行
1. 构建镜像：
   ```bash
   cd mcp
   docker build -t ech0-noise-mcp .
   ```
2. 以 Token 认证运行：
   ```bash
   docker run --rm -e NOTE_HOST=https://note.noisework.cn -e NOTE_TOKEN=你的_token ech0-noise-mcp
   ```
3. 以用户名密码运行：
   ```bash
   docker run --rm -e NOTE_HOST=https://note.noisework.cn ech0-noise-mcp
   # 在客户端调用“登录”工具设置 Cookie 会话
   ```

### docker-compose 一键启动（后端 + 前端静态 + MCP 同容器）
仓库根目录已有 `docker-compose.yml`。执行：
```bash
docker-compose up -d
```
- 服务 `my-app`：后端 Go + 前端静态，并内置 MCP 服务（Node）。
- 端口：应用 `1314`、MCP HTTP/SSE `1315`（皆映射至宿主）。

查看工具列表与流式调用：
```bash
curl http://localhost:1315/mcp/tools
curl -N -X POST http://localhost:1315/mcp/tool/搜索 -H 'Content-Type: application/json' -d '{"keyword":"#CDN"}'
```

### 镜像构建（多阶段）
- 带 MCP（同时提供 HTTP/SSE 与 Stdio）：
  ```bash
  docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --target final-mcp \
    -t noise233/echo-noise:latest \
    --push --no-cache .
  ```
- 不带 MCP（仅后端与静态前端）：
  ```bash
  docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --target final \
    -t noise233/echo-noise:latest \
    --push --no-cache .
  ```
- 如需单架构或自定义标签，可调整示例：
  ```bash
  docker buildx build --platform linux/amd64 --target final -t noise233/echo-noise:last --push --no-cache .
  ```

### HTTP 与 SSE
- 开启 HTTP 与 SSE：设置 `NOTE_HTTP_PORT` 启动服务端口
  ```bash
  NOTE_HOST=https://note.noisework.cn \
  NOTE_TOKEN=你的_token \
  NOTE_HTTP_PORT=1315 \
  npm start
  ```
- 列出工具：`GET /mcp/tools`
  ```bash
  curl http://localhost:1315/mcp/tools
  ```
- 流式调用工具：`POST /mcp/tool/{name}` 返回换行分隔 JSON
  ```bash
  curl -N -X POST http://localhost:1315/mcp/tool/搜索 \
    -H 'Content-Type: application/json' \
    -d '{"keyword":"#CDN","page":1,"pageSize":10}'
  ```
- SSE 订阅：`GET /mcp/sse`，推送 `tool_start`、`tool_end`
  ```bash
  curl -N http://localhost:1315/mcp/sse
  ```

#### 远程 URL 配置示例（不在本地运行进程）
- 一些客户端支持以 URL 订阅远程 SSE（无需本地启动 MCP 进程）。请使用纯净 URL 字符串，不要添加反引号或多余空格。
  ```json
  {
    "mcpServers": {
      "ech0-noise-mcp-http": {
        "url": "http://<你的服务器IP>:1315/mcp/sse"
      }
    }
  }
  ```
- 服务器端要求：容器以 `final-mcp` 目标运行，并设置 `NOTE_HTTP_PORT=1315`、映射端口 `-p 1315:1315`、`NOTE_HOST=http://localhost:1314`。
- 工具调用需通过 HTTP 端点：`POST http://<你的服务器IP>:1315/mcp/tool/{name}`，请求体为 JSON；SSE 用于事件订阅与握手展示。

#### SSE 工具执行（仅使用 SSE 触发并流式返回）
- 新增端点：`GET /mcp/sse/tool/{name}?input=<urlencoded-json>`
- 行为：立即触发工具执行，并以事件流返回：`tool_start`、`tool_data`、`tool_end`
- 示例：搜索
  ```bash
  curl -N "http://<服务器IP>:1315/mcp/sse/tool/search?input=%7B%22query%22:%22welcome%22,%22page%22:1,%22pageSize%22:10%7D"
  ```
- 示例：发布（容器端需已设置 `NOTE_TOKEN` 或先登录）
  ```bash
  curl -N "http://<服务器IP>:1315/mcp/sse/tool/publish?input=%7B%22type%22:%22text%22,%22content%22:%22测试内容%22%7D"
  ```
- 提示：`input` 为 URL 编码的 JSON；不建议在 URL 中使用反引号或空格。

### 本地连接远程 MCP（多种方式）

- 方式 A：SSE URL 订阅（不运行本地进程）
  - 适合仅订阅事件与握手展示，不进行完整 MCP 调用
  - 配置：`type: sse`，`url: http://<服务器IP>:1315/mcp/sse`
  - 注意：URL 必须是纯字符串，不要包含反引号或多余空格
  - 排错：
    - 确认服务器容器已设置 `NOTE_HTTP_PORT=1315` 并映射 `-p 1315:1315`
    - 服务器内执行 `curl -N http://localhost:1315/mcp/sse` 应返回事件流
    - 外网执行 `curl -N http://<服务器IP>:1315/mcp/sse` 验证链路；若失败检查防火墙/安全组开放 `1315/tcp`
    - 经由反向代理时需透传 `text/event-stream` 与 `Connection: keep-alive`

- 方式 B：Stdio 通过 docker exec 在容器内运行（本地/远程取决于 Docker 上下文）
  - 适合本地不安装 Node，通过容器完成握手与交互
  - 示例：
    ```json
    {
      "mcpServers": {
        "ech0-noise-mcp-stdio": {
          "command": "docker",
          "args": [
            "exec",
            "-i",
            "-e", "NOTE_HOST=http://<服务器IP>:1314",
            "-e", "NOTE_HTTP_PORT=0",
            "-e", "NOTE_TOKEN=<你的Token>",
            "Ech0-Noise",
            "node",
            "/app/mcp/server.bundle.mjs"
          ]
        }
      }
    }
    ```
  - 说明：
    - `NOTE_HTTP_PORT=0` 仅握手，不在容器内监听 HTTP；工具调用可走 HTTP 端点
    - `docker exec` 默认连接“本机 Docker 守护进程”。要操作云服务器上的容器，需先切换到“远程 Docker 上下文”。
  - 远程 Docker 上下文配置：
    ```bash
    # 创建远程上下文（通过 SSH）
    docker context create ech0-remote --docker "host=ssh://<user>@<服务器IP>"
    # 切换到远程上下文
    docker context use ech0-remote
    # 验证当前指向远程：应列出云服务器上的容器
    docker ps
    # 现在再执行 docker exec 即针对远程容器
    docker exec -it Ech0-Noise node /app/mcp/server.bundle.mjs
    ```
  - 常见错误：`Error response from daemon: No such container: Ech0-Noise`
    - 原因：当前 Docker 指向本机而非远程；或容器名不一致
    - 处理：切换到远程上下文，或确认容器名与 `docker-compose.yml` 的 `container_name: Ech0-Noise` 保持一致

- 方式 C：Stdio 在本地运行 Node，连接远程后端
  - 适合本地已安装 Node，直接运行仓库内 `mcp/server.js`
  - 推荐（跨平台）写法：使用环境变量字段而不是 `env` 命令
    ```json
    {
      "mcpServers": {
        "ech0-noise-local-stdio": {
          "command": "node",
          "args": ["/absolute/path/to/mcp/server.js"],
          "env": {
            "NOTE_HOST": "http://<服务器IP>:1314",
            "NOTE_HTTP_PORT": "0",
            "NOTE_TOKEN": "<你的Token>"
          }
        }
      }
    }
    ```
  - 可选（类 Unix）写法：使用 `env` 作为命令设置环境再运行 `node`
    ```json
    {
      "mcpServers": {
        "ech0-noise-local-stdio-env": {
          "command": "env",
          "args": [
            "NOTE_HOST=http://<服务器IP>:1314",
            "NOTE_HTTP_PORT=0",
            "NOTE_TOKEN=<你的Token>",
            "node",
            "/absolute/path/to/mcp/server.js"
          ]
        }
      }
    }
    ```
  - 提示：
    - 不要在 URL 中添加反引号或空格
    - 本地运行需提前 `cd mcp && npm install`
    - 为避免端口冲突，保持 `NOTE_HTTP_PORT=0`（仅 stdio 握手）；如需开启 HTTP/SSE，请设置为非 0 的端口并确认未被占用

  - 路径来源与依赖
    - `mcp/server.js` 位于本仓库根目录的 `mcp` 子目录中；请将示例中的绝对路径替换为你本机的实际克隆路径
      - macOS 示例：`/Library/Github/Ech0-Noise/mcp/server.js`
      - Windows 示例：`C:\Users\<你>\Github\Ech0-Noise\mcp\server.js`
      - Linux 示例：`/home/<你>/Github/Ech0-Noise/mcp/server.js`
    - 运行前在 `mcp` 目录安装依赖：`npm install`
    - Node.js 版本建议 `>= 20`（启用 ESM 与内置 `fetch`），否则可能出现模块或网络调用问题
    - 环境变量必须是纯字符串（不要加反引号或空格）：
      - 正确：`"NOTE_HOST": "http://154.219.122.129:1314"`
      - 错误：`"NOTE_HOST": " \`http://154.219.122.129:1314\` "`
    - 如果你不希望在本地安装 Node，可改用“方式 B：docker exec 远程容器”，命令中路径为容器内的 `/app/mcp/server.bundle.mjs`

- 方式 D：SSH + Stdio（无需本地 Docker）
  - 适合不在本地运行 Docker，通过 SSH 在远端启动 MCP 进程，STDIN/STDOUT 直接与本地客户端握手。
  - 远端要求：
    - 安装 Node.js `>= 20`
    - 已存在 MCP 文件：容器内单文件 `/app/mcp/server.bundle.mjs`（`final-mcp` 目标）或源码 `mcp/server.js`
    - 开放 SSH，建议使用非 `root` 用户（例如 `ubuntu`、`ec2-user`、`debian` 等）
  - 生成与安装 SSH 密钥（本地 macOS/Linux）：
    ```bash
    # 生成 Ed25519 密钥（推荐设置强口令；如不需要可回车留空）
    ssh-keygen -t ed25519 -a 100 -C "ech0" -f ~/.ssh/ech0_ed25519
    
    # 加载到本地代理（可选）
    ssh-add ~/.ssh/ech0_ed25519
    
    # 将公钥安装到远端用户（示例：ubuntu@<服务器IP>）
    ssh ubuntu@<服务器IP> 'mkdir -p ~/.ssh && chmod 700 ~/.ssh && cat >> ~/.ssh/authorized_keys' < ~/.ssh/ech0_ed25519.pub
    # 在远端设置权限
    # （如无法登录远端，可通过云平台控制台或已有账户手动写入 authorized_keys 并设置权限）
    chmod 600 ~/.ssh/authorized_keys
    
    # 测试 SSH 连通性
    ssh -i ~/.ssh/ech0_ed25519 ubuntu@<服务器IP> uname -a
    ```
  - MCP 配置示例（SSH + Stdio，纯字符串环境变量）：
    ```json
    {
      "mcpServers": {
        "ech0-mcp-stdio": {
          "command": "ssh",
          "args": [
            "-o", "IdentitiesOnly=yes",
            "-i", "/Users/<你>/.ssh/ech0_ed25519",
            "ubuntu@<服务器IP>",
            "bash", "-lc",
            "NOTE_HOST=http://<服务器IP>:1314 NOTE_HTTP_PORT=0 NOTE_TOKEN=<你的Token> node /app/mcp/server.bundle.mjs"
          ]
        }
      }
    }
    ```
    - 如远端 SSH 非 `22` 端口，加入 `-p <端口>` 参数。
    - 不要在 `NOTE_HOST` 中加入反引号或多余空格。
    - 若容器未内置 `server.bundle.mjs`，可改为路径 `node /absolute/path/to/mcp/server.js`（确保远端已安装依赖）。
  - 安全建议：
    - 为私钥设置口令（passphrase），避免泄露风险。
    - 尽量使用短期 Token；或在客户端先调用 `login` 工具获取会话后再进行敏感操作。
  - 常见排错：
    - `Permission denied (publickey,...)`：确认公钥已写入正确用户的 `~/.ssh/authorized_keys`，并设置权限 `700 ~/.ssh`、`600 ~/.ssh/authorized_keys`；确认使用的用户与 `-i` 指定的密钥匹配；很多云主机默认禁用 `root` 密钥登录。
    - 仅密码可用：移除 `BatchMode=yes` 才会提示密码，但不少 MCP 客户端不支持交互密码；建议改用密钥。
    - 环境变量格式：`NOTE_HOST` 必须是纯 URL 字符串；避免反引号与前后空格。
    - 远端 Node 版本过低：升级到 `>= 20` 以确保 ESM 与内置 `fetch` 可用。

#### 前端联动
- 置顶联动：前端列表会将置顶项排在顶部，取消置顶后按时间顺序排列（逻辑见 `web/components/index/MessageList.vue:260-264`）。
- 更新联动：修改内容后列表即时读取更新，无需重建；详情页 `消息` 工具可用于验证。
- 设置联动：通过 `settings` 更新后端配置，`配置` 工具可直接获取最新前端配置用于前端渲染。

### 本地开发同时启动 MCP（可选）
- 启动后端（Go）：
  ```bash
  go run ./cmd/server/main.go
  ```
- 启动前端（Nuxt 开发）：
  ```bash
  cd web && npm install && npm run dev
  ```
> 提示：容器集成版本已自动启动 MCP；本地开发时如需单独运行 MCP，可执行：
```bash
cd mcp && npm install && NOTE_HOST=http://localhost:1314 NOTE_HTTP_PORT=1315 npm start
```

### 桌面应用内自带 MCP（示例）
在 Electron 主进程中以子进程方式启动 MCP（避免用户二次启动）：
```js
// main.js（Electron）
import { spawn } from 'node:child_process'
const mcp = spawn('node', ['server.js'], { cwd: path.join(app.getAppPath(), 'mcp'), env: { NOTE_HOST: 'http://localhost:1314', NOTE_HTTP_PORT: '1315' } })
mcp.stdout.on('data', (d) => console.log('[mcp]', d.toString()))
mcp.stderr.on('data', (d) => console.error('[mcp]', d.toString()))
app.on('before-quit', () => { try { mcp.kill() } catch {} })
```
客户端可通过 Stdio（子进程）、或 `http://localhost:1315/mcp/tool/{name}` 流式调用 MCP 工具。

### MCP 客户端接入示例（Node）
使用 @modelcontextprotocol/sdk 连接并调用工具：
```js
import { Client } from '@modelcontextprotocol/sdk/client/index.js'
import { StdioClientTransport } from '@modelcontextprotocol/sdk/client/stdio.js'

const transport = new StdioClientTransport({
  command: 'node',
  args: ['mcp/server.js'],
  env: { NOTE_HOST: 'https://note.noisework.cn', NOTE_TOKEN: '你的_token' }
})

const client = new Client()
await client.connect(transport)

// 搜索
const res = await client.callTool({
  name: '搜索',
  input: { keyword: '#CDN', page: 1, pageSize: 10 }
})
console.log(res)

// 发布
await client.callTool({ name: '发布', input: { content: '测试内容', private: false } })

// 删除
await client.callTool({ name: '删除', input: { id: '123' } })
```

### 常用工具与入参
- 搜索/`search`：`{ keyword, page?, pageSize? }`（支持 `#标签`）
- 发布/`publish`：`{ content, private?, imageURL? }`
- 删除/`delete`：`{ id }`
- 更新/`update`：`{ id, content }`
- 消息/`消息`：`{ id }`
- 页面/`页面`：`{ page?, pageSize? }`
- 状态/`状态`：无入参
- 日历/`日历`：无入参
- 配置/`配置`：无入参
- 登录/`登录`：`{ username, password }`（设置 Cookie 会话）
- 令牌/`令牌`：无入参（基于会话生成新 Token）
- RSS/`RSS`：无入参（返回全文 XML）

### 与现有 API 的对应关系
- 搜索：`/api/messages/search` 或 `/api/messages/tags/:tag`
- 发布：`/api/messages`（会话）、`/api/token/messages`（Token）
- 删除/更新：`/api/messages/:id`
- 消息：`/api/messages/:id`
- 页面：`/api/messages/page`
- 状态：`/api/status`
- 日历：`/api/messages/calendar`
- 配置：`/api/frontend/config`
- 登录：`/api/login`（读取 `Set-Cookie`）
- 令牌：`/api/user/token/regenerate`
- RSS：`/rss`

### 接入其它客户端
- 绝大多数 MCP 客户端支持以 Stdio 方式接入：将命令设置为 `node mcp/server.js`，并传入环境变量 `NOTE_HOST`、`NOTE_TOKEN` 或在运行后调用“登录”。
- 工具名同时提供中文与英文，便于自然语言调用。

### 客户端配置示例

#### Claude Desktop

```json
{
  "mcpServers": {
    "ech0-noise-mcp-stdio": {
      "command": "docker",
      "args": [
        "exec",
        "-i",
        "-e", "NOTE_HOST=http://localhost:1314",
        "-e", "NOTE_HTTP_PORT=0",
        "-e", "NOTE_TOKEN=你的后台TOKEN",
        "Ech0-Noise",
        "node",
        "/app/mcp/server.bundle.mjs"
      ]
    },
    "ech0-noise-mcp-sse-monitor": {
      "type": "sse",
      "url": "http://localhost:1315/mcp/sse"
    }
  }
}
```

#### Trae IDE

```json
{
  "mcpServers": {
    "ech0-noise-mcp-stdio": {
      "command": "docker",
      "args": [
        "exec",
        "-i",
        "-e", "NOTE_HOST=http://localhost:1314",
        "-e", "NOTE_HTTP_PORT=0",
        "-e", "NOTE_TOKEN=你的后台TOKEN",
        "Ech0-Noise",
        "node",
        "/app/mcp/server.bundle.mjs"
      ]
    },
    "ech0-noise-mcp-sse-monitor": {
      "type": "sse",
      "url": "http://localhost:1315/mcp/sse"
    }
  }
}
```

#### Cherry Studio

```json
{
  "mcpServers": {
    "ech0-noise-mcp-sse": {
      "type": "sse",
      "url": "http://localhost:1315/mcp/sse"
    }
  }
}
```

#### Cursor

```json
{
  "mcpServers": {
    "ech0-noise-mcp-stdio": {
      "command": "docker",
      "args": [
        "exec",
        "-i",
        "-e", "NOTE_HOST=http://localhost:1314",
        "-e", "NOTE_HTTP_PORT=0",
        "-e", "NOTE_TOKEN=你的后台TOKEN",
        "Ech0-Noise",
        "node",
        "/app/mcp/server.bundle.mjs"
      ]
    },
    "ech0-noise-mcp-sse-monitor": {
      "type": "sse",
      "url": "http://localhost:1315/mcp/sse"
    }
  }
}
```

提示：

- `stdio` 用于握手与完整交互；为避免容器端口冲突，第二实例设置 `NOTE_HTTP_PORT=0`。
- `sse` 用于事件订阅与握手展示，不承担完整 MCP 交互。需要调用工具时使用 HTTP 或 `stdio`。
