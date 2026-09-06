FROM node:22-alpine AS builder

WORKDIR /app

RUN corepack enable

COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

COPY web/. .
COPY shared /shared

ARG APP_VERSION=dev
ARG BUILD_COMMIT=unknown
ENV APP_VERSION=${APP_VERSION} \
    BUILD_COMMIT=${BUILD_COMMIT}

RUN pnpm build

FROM node:22-alpine AS runtime

WORKDIR /app

ARG APP_VERSION=dev
ARG BUILD_COMMIT=unknown
ENV NODE_ENV=production \
    APP_VERSION=${APP_VERSION} \
    BUILD_COMMIT=${BUILD_COMMIT}

RUN apk add --no-cache su-exec \
    && addgroup -g 10001 -S app \
    && adduser -u 10001 -S app -G app \
    && corepack enable

COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile --prod

COPY --from=builder /app/build /app/build

# Install fonts for OG image rendering (resvg on musl cannot load woff2).
# Noto Serif SC (OTF) – CJK title & subtitle
# Google Sans Code (TTF) – tags & site name
ADD https://github.com/notofonts/noto-cjk/releases/download/Serif2.003/14_NotoSerifSC.zip /tmp/NotoSerifSC.zip
ADD https://github.com/googlefonts/googlesans-code/releases/download/v6.001/GoogleSansCode-v6.001.zip /tmp/GoogleSansCode.zip
RUN apk add --no-cache fontconfig && \
    mkdir -p /usr/share/fonts/og && \
    unzip -q -j /tmp/NotoSerifSC.zip 'SubsetOTF/SC/NotoSerifSC-Regular.otf' 'SubsetOTF/SC/NotoSerifSC-Bold.otf' -d /usr/share/fonts/og/ && \
    unzip -q -j /tmp/GoogleSansCode.zip 'variable/*.ttf' -d /usr/share/fonts/og/ && \
    fc-cache -f && \
    rm -f /tmp/NotoSerifSC.zip /tmp/GoogleSansCode.zip

COPY deploy/docker/renderer-entrypoint.sh /usr/local/bin/renderer-entrypoint.sh
RUN chmod +x /usr/local/bin/renderer-entrypoint.sh

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/renderer-entrypoint.sh"]
