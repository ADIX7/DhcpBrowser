# Backend
FROM golang:1.22.3-alpine AS backend_build
WORKDIR /app

COPY backend/go.mod ./
RUN go mod download

COPY backend/* ./

RUN go build -o /dhcpbrowser

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
FROM scratch AS final

COPY --from=backend_build /dhcpbrowser /dhcpbrowser
COPY --from=frontend_build /app/dist /static

EXPOSE 8090

CMD [ "/dhcpbrowser" ]
