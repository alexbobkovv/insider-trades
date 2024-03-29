/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    TRADES_GATEWAY_URL: 'http://localhost:8082/api-gateway/v1/trade-views',
  },
};

module.exports = nextConfig;
