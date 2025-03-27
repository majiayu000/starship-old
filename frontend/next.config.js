/** @type {import('next').NextConfig} */
const nextConfig = {
  // 添加代理配置
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*', // 后端API地址
      },
    ];
  },
};

module.exports = nextConfig; 