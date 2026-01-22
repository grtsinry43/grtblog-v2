export interface MarkdownComponentAttribute {
	key: string;
	label: string;
	placeholder?: string;
	defaultValue?: string;
	inputType?: 'text' | 'switch' | 'number';
}

export interface MarkdownComponentDefinition {
	name: string;
	label: string;
	description?: string;
	attrs: MarkdownComponentAttribute[];
	insertTemplate: string;
}

export const markdownComponents: MarkdownComponentDefinition[] = [
	{
		name: 'gallery',
		label: 'Gallery',
		description: '相册组件',
		attrs: [],
		insertTemplate: '::: gallery\n\n:::'
	},
	{
		name: 'callout',
		label: 'Callout',
		description: '提示框',
		attrs: [],
		insertTemplate: '::: callout\n\n:::'
	},
	{
		name: 'timeline',
		label: 'Timeline',
		description: '时间轴',
		attrs: [],
		insertTemplate: '::: timeline\n\n:::'
	},
	{
		name: 'year-card',
		label: 'Year Card',
		description: '年终总结卡片',
		attrs: [
			{ key: 'url', label: '链接', placeholder: 'https://example.com' },
			{ key: 'title', label: '标题', placeholder: '2025 年终总结' },
			{ key: 'type', label: '类型', defaultValue: 'page' },
			{ key: 'cover', label: '封面图', placeholder: 'https://example.com/cover.jpg' },
			{ key: 'blur', label: '模糊度', defaultValue: '7px' }
		],
		insertTemplate:
			'::: year-card{url="" title="" type="page" cover="" blur="7px"}\n\n:::'
	},
	{
		name: 'link-card',
		label: 'Link Card',
		description: '链接卡片',
		attrs: [
			{ key: 'href', label: '链接', placeholder: '/path/to/page' },
			{ key: 'title', label: '标题', placeholder: '标题' },
			{ key: 'desc', label: '描述', placeholder: '描述' },
			{ key: 'newtab', label: '新窗口', defaultValue: 'true', inputType: 'switch' }
		],
		insertTemplate: '::: link-card{href="" title="" desc="" newtab="true"}\n\n:::'
	}
];

export const markdownComponentNames = new Set(markdownComponents.map((component) => component.name));

export const getMarkdownComponent = (name?: string) =>
	markdownComponents.find((component) => component.name === name);

export interface ParsedComponentInfo {
	name: string;
	attrs: Record<string, string>;
	rawAttrs: string;
}

const attrRegex = /([A-Za-z][\w-]*)\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s}]+))/g;

export const parseComponentAttributes = (raw: string) => {
	if (!raw) {
		return {};
	}
	const trimmed = raw.trim();
	const content = trimmed.startsWith('{') ? trimmed.slice(1) : trimmed;
	const normalized = content.endsWith('}') ? content.slice(0, -1) : content;
	const attrs: Record<string, string> = {};
	let match: RegExpExecArray | null = null;

	while ((match = attrRegex.exec(normalized)) !== null) {
		const key = match[1];
		const value = match[2] ?? match[3] ?? match[4] ?? '';
		if (typeof key === 'string') {
			attrs[key] = value;
		}
	}

	return attrs;
};

export const parseComponentInfo = (info: string): ParsedComponentInfo => {
	const trimmed = info.trim();
	if (!trimmed) {
		return { name: 'unknown', attrs: {}, rawAttrs: '' };
	}

	let rest = trimmed;
	if (trimmed.startsWith('component')) {
		const match = trimmed.match(/^component\s+(.+)$/);
		rest = (match?.[1] || '').trim();
	}

	const match = rest.match(/^([^\s{]+)\s*(\{.*)?$/);
	const name = (match?.[1] || 'unknown').trim();
	const rawAttrs = match?.[2]?.trim() || '';

	return { name, attrs: parseComponentAttributes(rawAttrs), rawAttrs };
};
