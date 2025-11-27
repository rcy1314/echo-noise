# 构建阶段：前端
FROM node:22.14.0-alpine AS frontend-build

# 设置工作目录
WORKDIR /app/web

# 复制前端依赖文件并安装依赖
COPY ./web/package.json ./web/package-lock.json* ./
# 使用稳定镜像源以避免 EAI_AGAIN
RUN npm config set registry https://registry.npmmirror.com && \
    npm ci --omit=dev --registry=https://registry.npmmirror.com

# 更新 Browserslist 数据，避免 caniuse-lite 过期警告
RUN npx --yes update-browserslist-db@latest

# 复制前端源代码并构建
COPY ./web/ .
RUN npm run generate

# 将构建结果复制到公共目录
RUN mkdir -p /app/public && cp -r .output/public/* /app/public/

# 构建阶段：后端
FROM golang:1.24.1-alpine AS backend-build

# 设置环境变量
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

# 设置工作目录
WORKDIR /app

# 配置 APK 镜像源并安装构建依赖
RUN echo "https://mirrors.aliyun.com/alpine/v3.21/main" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.21/community" >> /etc/apk/repositories && \
    apk update && apk add --no-cache build-base

# 复制 Go 模块文件并下载依赖
COPY ./go.mod ./go.sum ./
RUN go mod download

# 复制项目文件
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./config ./config

# 整理依赖，生成完整 go.sum（修复禁用 CGO 后构建缺少依赖的问题）
RUN go mod tidy

# 创建必要的目录并设置权限
RUN mkdir -p /app/data /app/public && chmod -R 755 /app/data

# 编译 Go 应用
RUN go build -trimpath -ldflags "-s -w -buildid=" -o /app/noise ./cmd/server/main.go

# MCP 构建阶段（打包为单文件，避免在最终镜像中保留 node_modules）
FROM node:20-alpine AS mcp-build
WORKDIR /app/mcp
COPY ./mcp/package.json ./
RUN npm config set registry https://registry.npmmirror.com && \
    npm install --omit=dev --registry=https://registry.npmmirror.com
COPY ./mcp/server.js ./server.js
RUN npx --yes esbuild@0.23.0 server.js \
    --bundle \
    --platform=node \
    --format=esm \
    --outfile=server.bundle.mjs

# 运行时阶段
FROM alpine:3.21 AS final

# 可选：是否使用 UPX 压缩二进制（1=启用，0=禁用）
ARG USE_UPX=1

# 设置工作目录
WORKDIR /app

# 从后端构建阶段复制配置文件和二进制文件
COPY --from=backend-build /app/config /app/config
COPY --from=backend-build /app/noise /app/noise

# 从前端构建阶段复制静态文件
COPY --from=frontend-build /app/public /app/public

# 复制 MCP 打包后的单文件（不包含 node_modules，减小镜像体积）
 

# 按需裁剪静态字体：保留仅在 CSS 中引用的 woff2
RUN set -eux; \
    if [ -d /app/public/_nuxt ]; then \
      cd /app/public/_nuxt; \
      ls *.woff2 >/dev/null 2>&1 || exit 0; \
      grep -Eoh '[A-Za-z0-9._-]+\.woff2' -- *.css 2>/dev/null | sort -u > keep.list || true; \
      for f in *.woff2; do \
        grep -qx "$f" keep.list || rm -f "$f"; \
      done; \
      rm -f keep.list; \
    fi

# 额外剔除：删除前端 sourcemap 与许可证文件，减少镜像体积
ARG KEEP_SOURCEMAPS=0
RUN set -eux; \
    if [ "$KEEP_SOURCEMAPS" != "1" ]; then \
      find /app/public -type f -name '*.map' -delete; \
      find /app/public -type f -name '*LICENSE*' -delete; \
      find /app/public -type f -name '.DS_Store' -delete; \
    fi


# 更换 Alpine 镜像源
RUN echo "https://mirrors.aliyun.com/alpine/v3.21/main" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.21/community" >> /etc/apk/repositories

# 安装运行时所需的工具
RUN apk update && \
    apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/*

# 创建数据和图片目录
RUN mkdir -p /app/data/images && \
    chmod -R 755 /app/data

# 内置 SQLite 数据库（初始数据），以便首次启动有内容
COPY ./data/noise.db /app/data/

# 可选：使用 UPX 压缩二进制以减小体积（默认启用，影响极小）
RUN if [ "$USE_UPX" = "1" ]; then \
      apk add --no-cache upx; \
      upx --best --lzma -q /app/noise || true; \
    fi

# 暴露应用端口
EXPOSE 1314
EXPOSE 1315

# 启动后端与 MCP（MCP 后台运行，Go 服务为主进程）
CMD ["/app/noise"]

FROM final AS final-mcp
RUN echo "https://mirrors.aliyun.com/alpine/v3.21/main" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.21/community" >> /etc/apk/repositories && \
    apk update && \
    apk add --no-cache nodejs && \
    rm -rf /var/cache/apk/*
COPY --from=mcp-build /app/mcp/server.bundle.mjs /app/mcp/server.bundle.mjs
ENV NOTE_HTTP_PORT=1315
EXPOSE 1314
EXPOSE 1315
CMD ["sh","-c","node /app/mcp/server.bundle.mjs & exec /app/noise"]

FROM node:20-alpine AS mcp-final
WORKDIR /app/mcp
COPY --from=mcp-build /app/mcp/server.bundle.mjs /app/mcp/server.bundle.mjs
ENV NOTE_HTTP_PORT=1315
EXPOSE 1315
CMD ["node","/app/mcp/server.bundle.mjs"]
