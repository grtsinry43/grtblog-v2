<script module>
    import {defineMeta} from '@storybook/addon-svelte-csf';

    import Button from '../lib/components/ui/button/Button.svelte';

    const {Story} = defineMeta({
        component: Button,
        args: {
            label: 'Button',
            variant: 'primary',
            size: 'md',
            fullWidth: false,
            loading: false,
            disabled: false
        },
        argTypes: {
            label: {control: 'text'},
            variant: {control: 'select', options: ['primary', 'secondary', 'ghost', 'icon']},
            size: {control: 'select', options: ['sm', 'md', 'lg']},
            fullWidth: {control: 'boolean'},
            loading: {control: 'boolean'},
            disabled: {control: 'boolean'}
        }
    });
</script>

<Story name="Primary" let:args>
    <Button {...args} content={args.label} />
</Story>
