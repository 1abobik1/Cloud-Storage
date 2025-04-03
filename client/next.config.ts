import type { NextConfig } from "next";

const nextConfig: NextConfig = {

    async rewrites() {
      return [
        {
          source: "/api/:path*",  // Перехватываем запросы на /api/*
          destination: "http://localhost:8080/:path*", // Проксируем на бэкенд
        },
      ];
    },
  }

export default nextConfig;
