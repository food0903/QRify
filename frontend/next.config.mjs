/** @type {import('next').NextConfig} */
const nextConfig = {
    typescript: {
      ignoreBuildErrors: true,
    },
    eslint: {
      ignoreDuringBuilds: true,
    },
    async rewrites() {
      if (process.env.NODE_ENV === "development") {
        return [
          {
            source: '/api/:path*',
            destination: 'http://localhost:8080/:path*',
          },
        ];
      }
      return [];
    },
};

export default nextConfig