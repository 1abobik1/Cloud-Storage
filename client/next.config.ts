import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'https://api.escuelajs.co/api/v1/:path*' // Proxy to Backend
      }
    ]
  }
};

export default nextConfig;
