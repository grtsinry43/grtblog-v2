import { brand } from '$lib/shared/brand';

export const consoleLogInfo = (): void => {
    const fontMono = 'font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;';

    const nameStyle =
        `background: #27272a; color: #fff; padding: 3px 6px; border-radius: 0; font-weight: 700; font-size: 12px; ${fontMono}`;

    const verStyle = 
        `background: #52525b; color: #fff; padding: 3px 6px; border-radius: 0; font-weight: 500; font-size: 12px; ${fontMono}`;

    const sloganStyle =
        `font-weight: 600; font-size: 10px; padding-left: 6px; ${fontMono} opacity: 0.9; margin-top: 6px;`;

    const authorStyle =
        `font-weight: 700; font-size: 12px; ${fontMono} margin-top: 4px; display: block;`;

    const linkStyle =
        `font-weight: 400; font-size: 11px; text-decoration: underline; opacity: 0.7; padding-left: 4px; ${fontMono} cursor: pointer;`;

    const mottoStyle =
        `font-style: italic; font-size: 11px; opacity: 0.6; padding-top: 4px; line-height: 1.5; font-family: system-ui, sans-serif;`;

    console.log(
        `%c${brand.name}%c${brand.version}%c ${brand.slogan}`,
        nameStyle,
        verStyle,
        sloganStyle
    );

    console.log(
        `by %c@${brand.author}%c${brand.github}\n%c${brand.motto}`,
        authorStyle,
        linkStyle,
        mottoStyle
    );
};