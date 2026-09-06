-- +goose Up
INSERT INTO sys_config (config_key, value, is_sensitive, group_path, label, description, value_type, sort, meta)
VALUES
('site.discovery.intro', '', false, 'site/discovery', '博客阅读简介', '用于 llms.txt 的简短博客介绍：写作主题、内容定位、适合哪些读者。留空使用网站描述。此处内容会公开。', 'string', 10, '{"inputType":"textarea"}'::jsonb),
('site.discovery.author', '', false, 'site/discovery', '作者背景', '公开介绍作者的经历、兴趣与写作视角。请勿填写私人联系方式或其他非公开信息。', 'string', 20, '{"inputType":"textarea"}'::jsonb),
('site.discovery.guidance', '', false, 'site/discovery', '阅读与引用说明', '补充术语解释、内容时效性、署名与引用偏好。支持 Markdown；默认指南已经说明原文引用、时间与外部来源的区别。', 'string', 30, '{"inputType":"textarea"}'::jsonb),
('site.discovery.featured_paths', '', false, 'site/discovery', '精选阅读路径', '每行一个站内页面路径，如 /posts/hello/ 或 /about/。只展示存在且公开的文章、手记或独立页面；撤稿后自动移除。', 'string', 40, '{"inputType":"textarea","placeholder":"/about/\n/posts/hello/"}'::jsonb),
('site.discovery.robots', E'User-agent: *\nDisallow: /admin', false, 'site/discovery', 'robots.txt 爬取策略', '按 robots.txt 格式填写，可为不同 User-agent 设置 Allow / Disallow。默认及留空时禁止抓取 /admin，其余路径允许抓取。系统自动附加本站 Sitemap 地址，无需手填；后台页面另外固定声明 noindex。', 'string', 50, '{"inputType":"textarea","placeholder":"User-agent: *\nDisallow: /admin"}'::jsonb)
ON CONFLICT (config_key) DO NOTHING;

-- 为已配置页脚分组的站点追加「如果你是 Agent」导航分组
WITH parsed AS (
    SELECT
        config_key,
        CASE
            WHEN value IS NULL OR btrim(value) = '' THEN '{}'::jsonb
            WHEN jsonb_typeof(value::jsonb) = 'object' THEN value::jsonb
            ELSE '{}'::jsonb
        END AS theme
    FROM sys_config
    WHERE config_key = 'site.theme_extend_info'
),
footer_layer AS (
    SELECT
        config_key,
        theme,
        CASE
            WHEN jsonb_typeof(theme -> 'footer') = 'object' THEN theme -> 'footer'
            ELSE '{}'::jsonb
        END AS footer
    FROM parsed
),
patched AS (
    SELECT
        config_key,
        jsonb_set(
            theme,
            '{footer}',
            footer || jsonb_build_object(
                'sections',
                CASE
                    WHEN jsonb_typeof(footer -> 'sections') = 'array'
                        AND EXISTS (
                            SELECT 1
                            FROM jsonb_array_elements(footer -> 'sections') AS s
                            WHERE s ->> 'title' = '如果你是 Agent >'
                        )
                    THEN footer -> 'sections'
                    WHEN jsonb_typeof(footer -> 'sections') = 'array'
                    THEN footer -> 'sections' || '[
                        {"title": "如果你是 Agent >", "links": [{"name": "站点地图", "href": "/sitemap"}, {"name": "AI 阅读指南", "href": "/llms.txt"}]}
                    ]'::jsonb
                    ELSE footer -> 'sections'
                END
            ),
            true
        ) AS value_json
    FROM footer_layer
)
UPDATE sys_config AS sc
SET value = patched.value_json::text,
    updated_at = now()
FROM patched
WHERE sc.config_key = patched.config_key
    AND jsonb_typeof(patched.value_json -> 'footer' -> 'sections') = 'array';

-- +goose Down
WITH parsed AS (
    SELECT
        config_key,
        CASE
            WHEN value IS NULL OR btrim(value) = '' THEN '{}'::jsonb
            WHEN jsonb_typeof(value::jsonb) = 'object' THEN value::jsonb
            ELSE '{}'::jsonb
        END AS theme
    FROM sys_config
    WHERE config_key = 'site.theme_extend_info'
),
footer_layer AS (
    SELECT
        config_key,
        theme,
        CASE
            WHEN jsonb_typeof(theme -> 'footer') = 'object' THEN theme -> 'footer'
            ELSE '{}'::jsonb
        END AS footer
    FROM parsed
),
sections_layer AS (
    SELECT
        config_key,
        theme,
        footer,
        CASE
            WHEN jsonb_typeof(footer -> 'sections') = 'array' THEN footer -> 'sections'
            ELSE '[]'::jsonb
        END AS sections
    FROM footer_layer
),
stripped AS (
    SELECT
        config_key,
        theme,
        footer,
        COALESCE(
            (
                SELECT jsonb_agg(sec)
                FROM jsonb_array_elements(sections) AS sec
                WHERE sec ->> 'title' <> '如果你是 Agent >'
            ),
            '[]'::jsonb
        ) AS sections_stripped
    FROM sections_layer
),
patched AS (
    SELECT
        config_key,
        jsonb_set(theme, '{footer,sections}', sections_stripped, true) AS value_json
    FROM stripped
)
UPDATE sys_config AS sc
SET value = patched.value_json::text,
    updated_at = now()
FROM patched
WHERE sc.config_key = patched.config_key
    AND jsonb_typeof(sc.value::jsonb -> 'footer' -> 'sections') = 'array';

DELETE FROM sys_config WHERE config_key IN ('site.discovery.intro', 'site.discovery.author', 'site.discovery.guidance', 'site.discovery.featured_paths', 'site.discovery.robots');
