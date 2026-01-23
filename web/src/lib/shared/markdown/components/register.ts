import YearCard from './YearCard.svelte';
import LinkCard from './LinkCard.svelte';
import { registerMarkdownComponent } from './index';

registerMarkdownComponent('year-card', YearCard);
registerMarkdownComponent('link-card', LinkCard);
