import type { NextConfig } from "next";

const nextConfig: NextConfig = {
 
  
  
  async rewrites() {
    return [
      
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/:path*',
      },
      {
        source: '/api/:path*',
        destination: 'http://localhost:8081/:path*',
      },
    ];
  },
};

export default nextConfig;
