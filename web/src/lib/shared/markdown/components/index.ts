import { mount, unmount } from 'svelte';
import type { Component } from 'svelte';

type ComponentConstructor = Component<Record<string, any>>;

const registry = new Map<string, ComponentConstructor>();

const escapeHtml = (input: string) =>
	input.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');

const parseProps = (raw?: string) => {
	if (!raw) {
		return {};
	}
	try {
		const parsed = JSON.parse(raw);
		return typeof parsed === 'object' && parsed ? parsed : {};
	} catch {
		return {};
	}
};

const renderUnsupported = (el: HTMLElement, name: string, props: Record<string, any>) => {
	el.innerHTML = `
		<div class="md-component-fallback">
			<span class="md-component-fallback__label">组件暂不支持</span>
		</div>
	`;
};

export const registerMarkdownComponent = (name: string, component: ComponentConstructor) => {
	registry.set(name, component);
};

export const unregisterMarkdownComponent = (name: string) => {
	registry.delete(name);
};

export const mountMarkdownComponents = (root: HTMLElement) => {
	const placeholders = Array.from(root.querySelectorAll<HTMLElement>('.md-component-placeholder'));
	const instances: Array<unknown> = [];

	for (const el of placeholders) {
		const name = el.dataset.component?.trim() || '';
		const Component = registry.get(name);
		const props = parseProps(el.dataset.props);
		const contentHtml = el.innerHTML;

		if (!Component) {
			renderUnsupported(el, name, props);
			continue;
		}

		el.innerHTML = '';
		const instance = mount(Component, { target: el, props: { ...props, contentHtml } });
		instances.push(instance);
	}

	return () => {
		for (const instance of instances) {
			unmount(instance as never);
		}
	};
};
