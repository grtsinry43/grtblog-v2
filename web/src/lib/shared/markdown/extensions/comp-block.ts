import container from 'markdown-it-container';

import { markdownComponents, parseComponentInfo } from '../shared/components';

import type { MarkdownExtension } from '../types';

/**
 * Custom component block syntax.
 * Syntax: ::: component <name> :::
 * Renders a placeholder div with data-component and data-props attributes.
 */
export const componentBlockExtension: MarkdownExtension = (md) => {
	const registerContainer = (
		containerName: string,
		resolveComponentInfo: (tokenInfo: string) => ReturnType<typeof parseComponentInfo>,
		validate?: (params: string) => boolean
	) => {
		md.use(container as any, containerName, {
			validate,
			render: (tokens: any[], idx: number) => {
				if (tokens[idx].nesting === 1) {
					const { name, attrs } = resolveComponentInfo(tokens[idx].info);
					const propsJson = JSON.stringify(attrs);
					const propsAttr = propsJson !== '{}' ? ` data-props="${md.utils.escapeHtml(propsJson)}"` : '';
					return `<div class="md-component-placeholder" data-component="${md.utils.escapeHtml(name)}"${propsAttr}>`;
				}
				return '</div>\n';
			}
		});
	};

	registerContainer(
		'component',
		(info) => parseComponentInfo(info),
		(params) => /^component\s+/.test(params.trim())
	);

	markdownComponents.forEach((component) => {
		const prefix = component.name;
		registerContainer(
			component.name,
			(info) => parseComponentInfo(info),
			(params) => params.trim().startsWith(prefix)
		);
	});
};
