import type { Handle, HandleServerError } from '@sveltejs/kit';
import { ISR_DEPS_HEADER } from '$lib/server/isr-deps';

const STATIC_FALLBACK_HEADER = 'x-grt-static-fallback';
const STATIC_MISS_WARN_TTL_MS = 5 * 60 * 1000;
const BACKFILL_THROTTLE_TTL_MS = 10 * 1000;

const createThrottle = (ttlMs: number) => {
	const seenAt = new Map<string, number>();
	return (key: string, now: number): boolean => {
		const last = seenAt.get(key) ?? 0;
		if (now - last < ttlMs) {
			return false;
		}
		seenAt.set(key, now);

		// Keep the in-memory throttle map bounded for long-lived renderer processes.
		if (seenAt.size > 1000) {
			for (const [candidate, timestamp] of seenAt) {
				if (now - timestamp >= ttlMs) {
					seenAt.delete(candidate);
				}
			}
		}
		return true;
	};
};

const shouldWarnStaticMiss = createThrottle(STATIC_MISS_WARN_TTL_MS);
const shouldBackfill = createThrottle(BACKFILL_THROTTLE_TTL_MS);

const defaultInternalServerBaseURL = 'http://localhost:8080';

const resolveInternalServerBaseURL = (): string => {
	if (typeof process === 'undefined' || !process.env) return defaultInternalServerBaseURL;
	const raw = (process.env.INTERNAL_API_BASE_URL || '').trim();
	if (!raw) return defaultInternalServerBaseURL;
	// Strip /api/v2 suffix if present to get the server root.
	return raw.replace(/\/api\/v2\/?$/, '').replace(/\/+$/, '') || defaultInternalServerBaseURL;
};

/**
 * Map a fallback request pathname to the page path the ISR queue should
 * re-render, or null when the request is not a snapshot-able page
 * (client assets, API calls, image/data files, ...).
 */
const resolveBackfillPath = (pathname: string): string | null => {
	let candidate = pathname;
	if (candidate.endsWith('/__data.json')) {
		candidate = candidate.slice(0, -'/__data.json'.length) || '/';
	}
	if (candidate !== '/' && candidate.endsWith('/')) {
		candidate = candidate.replace(/\/+$/, '') || '/';
	}
	if (candidate.startsWith('/_app/') || candidate.startsWith('/api/')) return null;
	// Discovery stays dynamic until its publication/withdrawal lifecycle joins ISR.
	if (candidate === '/sitemap') return null;
	// Paths with a file extension are assets (og-image.png, favicon, ...), not pages.
	if (/\.[a-zA-Z0-9]+$/.test(candidate)) return null;
	return candidate;
};

/**
 * Static-first self-healing: when nginx misses a static snapshot and this SSR
 * fallback successfully renders the page, ask the server to regenerate the
 * snapshot so subsequent requests are served statically again.
 */
const backfillStaticSnapshot = (urlPath: string): void => {
	const endpoint = `${resolveInternalServerBaseURL()}/internal/isr/revalidate`;
	void fetch(endpoint, {
		method: 'POST',
		headers: { 'content-type': 'application/json' },
		body: JSON.stringify({ urlPath })
	}).catch((err) => {
		console.warn(`[renderer][isr-backfill] enqueue failed path=${urlPath} err=${err}`);
	});
};

const logServerResponse = (
	method: string,
	pathname: string,
	status: number,
	staticFallback: boolean,
	deps: string[]
) => {
	if (status < 400) return;
	const level = status >= 500 ? 'error' : 'warn';
	const extra = staticFallback ? ' staticFallback=1' : '';
	const depInfo = deps.length > 0 ? ` deps=${deps.join(',')}` : '';
	console[level](
		`[renderer][server-response] side=server code=${status} method=${method} path=${pathname}${extra}${depInfo}`
	);
};

export const handleError: HandleServerError = ({ error, event, status, message }) => {
	const detail =
		error instanceof Error
			? `${error.name}: ${error.message}${error.stack ? `\n${error.stack}` : ''}`
			: String(error);
	console.error(
		`[renderer][server-exception] side=server code=${status ?? 500} method=${event.request.method} path=${event.url.pathname} message=${message}\n${detail}`
	);
};

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.isrDeps = new Set<string>();

	const response = await resolve(event);
	const staticFallback = event.request.headers.get(STATIC_FALLBACK_HEADER) === '1';
	const depList = Array.from(event.locals.isrDeps).sort();
	logServerResponse(
		event.request.method,
		event.url.pathname,
		response.status,
		staticFallback,
		depList
	);
	if (staticFallback && event.locals.isrDeps.size > 0) {
		const now = Date.now();
		const warnKey = `${event.request.method}:${event.url.pathname}:${response.status}:${depList.join(',')}`;
		if (shouldWarnStaticMiss(warnKey, now)) {
			console.warn(
				`[renderer][isr-static-miss] ${event.request.method} ${event.url.pathname} status=${response.status} deps=${depList.join(',')}`
			);
		}
	}
	if (staticFallback && event.request.method === 'GET' && response.status < 400) {
		const backfillPath = resolveBackfillPath(event.url.pathname);
		if (backfillPath && shouldBackfill(backfillPath, Date.now())) {
			backfillStaticSnapshot(backfillPath);
		}
	}
	if (event.locals.isrDeps.size === 0) {
		return response;
	}

	const headers = new Headers(response.headers);
	headers.set(ISR_DEPS_HEADER, JSON.stringify(depList));
	return new Response(response.body, {
		status: response.status,
		statusText: response.statusText,
		headers
	});
};
