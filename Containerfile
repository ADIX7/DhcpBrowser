# Backend
FROM golang:1.22.3-alpine AS backend_build
WORKDIR /app

COPY backend/go.mod ./
RUN go mod download

COPY backend/* ./

RUN go build -o /dhcpbrowser

FROM alpine:3.13 as backend

COPY --from=backend_build /dhcpbrowser /dhcpbrowser

EXPOSE 8090

CMD [ "/dhcpbrowser" ]

# Frontend
FROM node:20-alpine AS frontend_build
WORKDIR /app

COPY frontend/package.json ./
COPY frontend/pnpm-lock.yaml ./
RUN corepack enable pnpm
# RUN corepack prepare pnpm@latest --activate
RUN pnpm install

COPY frontend .
RUN ls
RUN pnpm build

# Final
FROM backend as final
COPY --from=frontend_build /app/dist /static
