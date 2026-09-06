FROM node:22-alpine AS builder

WORKDIR /app

RUN corepack enable

COPY admin/package.json admin/pnpm-lock.yaml admin/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

COPY admin/. ./
COPY shared /shared

ARG VITE_APP_BASE=/admin/
ARG VITE_APP_NAME=Grtblog Admin
ARG VITE_APP_TITLE=管理后台
ARG VITE_WATERMARK_CONTENT=
ARG VITE_API_BASE_URL=/api/v2

ENV VITE_APP_BASE=${VITE_APP_BASE} \
    VITE_APP_NAME=${VITE_APP_NAME} \
    VITE_APP_TITLE=${VITE_APP_TITLE} \
    VITE_WATERMARK_CONTENT=${VITE_WATERMARK_CONTENT} \
    VITE_API_BASE_URL=${VITE_API_BASE_URL}

RUN pnpm build

FROM nginx:1.27-alpine AS runtime

COPY --from=builder /app/dist /usr/share/nginx/html/admin

RUN cat > /etc/nginx/conf.d/default.conf <<'NGINX'
server {
    listen 80;
    server_name _;

    root /usr/share/nginx/html;

    location = / {
        return 302 /admin/;
    }

    location /admin/ {
        try_files $uri $uri/ /admin/index.html;
    }
}
NGINX

EXPOSE 80
