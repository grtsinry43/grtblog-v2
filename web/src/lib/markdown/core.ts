import MarkdownIt from 'markdown-it';

import { componentBlockExtension } from './extensions/comp-block';

import type { MarkdownConfig } from './types';

const applyHeadingAnchorRule = (md: MarkdownIt) => {
	const defaultHeadingOpen =
		md.renderer.rules.heading_open ||
		((tokens, idx, options, _env, self) => self.renderToken(tokens, idx, options));

	md.renderer.rules.heading_open = (tokens, idx, options, env, self) => {
		const anchors = (env as any)?.headingAnchors as string[] | undefined;
		if (anchors && typeof (env as any)?.headingAnchorIndex === 'number') {
			const anchor = anchors[(env as any).headingAnchorIndex];
			if (anchor) {
				tokens[idx].attrSet('id', anchor);
				(env as any).headingAnchorIndex += 1;
			}
		}
		return defaultHeadingOpen(tokens, idx, options, env, self);
	};
};

export const createMarkdownIt = (config: MarkdownConfig = {}) => {
	const md = new MarkdownIt({
		html: true,
		linkify: true,
		typographer: true,
		...config.options
	});

	componentBlockExtension(md);
	config.extensions?.forEach((extension) => extension(md));
	applyHeadingAnchorRule(md);

	return md;
};
