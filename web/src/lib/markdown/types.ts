import type MarkdownIt from 'markdown-it';
import type { Options } from 'markdown-it';

export type MarkdownExtension = (md: MarkdownIt, options?: any) => void;

export type MarkdownConfig = {
	options?: Options;
	extensions?: MarkdownExtension[];
};
