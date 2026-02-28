import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Forward /api/* calls to the Go backend in dev
  async rewrites() {
    return process.env.NODE_ENV === "development"
      ? [
          {
            source: "/api/:path*",
            destination: "http://localhost:9999/api/:path*",
          },
        ]
      : [];
  },
};

export default nextConfig;
