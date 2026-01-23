import { createMarkdownIt } from '$lib/shared/markdown/core';

const markdown = createMarkdownIt();

export const renderMarkdown = (input: string, headingAnchors: string[] = []) =>
	markdown.render(input ?? '', { headingAnchors, headingAnchorIndex: 0 });
