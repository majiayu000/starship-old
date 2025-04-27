FROM node:18-alpine AS base

# 安装依赖
FROM base AS deps
WORKDIR /app

# 复制 package.json 和 package-lock.json
COPY package.json package-lock.json* ./

# 安装依赖
RUN npm ci

# 构建应用
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# 设置环境变量
ENV NEXT_TELEMETRY_DISABLED 1
ENV NODE_ENV production

# 构建应用
RUN npm run build

# 生产环境
FROM base AS runner
WORKDIR /app

# 设置环境变量
ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

# 创建非 root 用户
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# 复制构建产物
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# 切换到非 root 用户
USER nextjs

# 暴露端口
EXPOSE 5555

# 设置环境变量
ENV PORT 5555
ENV HOSTNAME "0.0.0.0"

# 启动应用
CMD ["node", "server.js"]
