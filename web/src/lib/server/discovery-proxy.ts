import type { RequestEvent } from '@sveltejs/kit';

/** Raw resources are proxied without the JSON API envelope or user credentials. */
export async function proxyDiscovery(event: RequestEvent): Promise<Response> {
	const configured = process.env.INTERNAL_API_BASE_URL?.trim() || 'http://localhost:8080/api/v2';
	const base = configured.replace(/\/api\/v2\/?$/, '').replace(/\/+$/, '');
	const upstream = new URL(`${base}/api/v2/public/discovery/resource${event.url.pathname}`);
	upstream.search = event.url.search;
	const headers = new Headers();
	const etag = event.request.headers.get('if-none-match');
	if (etag) headers.set('if-none-match', etag);
	try {
		const response = await event.fetch(upstream, {
			method: event.request.method,
			headers,
			signal: AbortSignal.timeout(15000),
			redirect: 'manual'
		});
		const output = new Headers();
		for (const name of [
			'content-type',
			'cache-control',
			'etag',
			'link',
			'x-robots-tag',
			'x-content-type-options'
		]) {
			const value = response.headers.get(name);
			if (value) output.set(name, value);
		}
		return new Response(event.request.method === 'HEAD' ? null : response.body, {
			status: response.status,
			headers: output
		});
	} catch (error) {
		console.error('[discovery] upstream unavailable', error);
		return new Response(
			event.request.method === 'HEAD' ? null : 'Discovery temporarily unavailable',
			{
				status: 503,
				headers: { 'cache-control': 'no-store', 'content-type': 'text/plain; charset=utf-8' }
			}
		);
	}
}
