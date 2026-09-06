import { proxyDiscovery } from '$lib/server/discovery-proxy';

export const prerender = false;
export const trailingSlash = 'never';
export const GET = proxyDiscovery;
export const HEAD = proxyDiscovery;
