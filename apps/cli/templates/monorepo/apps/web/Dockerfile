FROM oven/bun:1.2.13 AS builder

WORKDIR /repo

ARG VITE_API_URL=http://localhost:8080
ARG VITE_ENV=development
ARG VITE_GOOGLE_AUTH_ENABLED=false

ENV VITE_API_URL=${VITE_API_URL}
ENV VITE_ENV=${VITE_ENV}
ENV VITE_GOOGLE_AUTH_ENABLED=${VITE_GOOGLE_AUTH_ENABLED}

COPY package.json bun.lock turbo.json ./
COPY packages ./packages
COPY apps/api/package.json ./apps/api/
COPY apps/web ./apps/web

RUN bun install --frozen-lockfile --ignore-scripts
RUN bun run build --filter=./apps/web

FROM nginx:1.27-alpine

COPY apps/web/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /repo/apps/web/dist /usr/share/nginx/html

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
