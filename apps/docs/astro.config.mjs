// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'cliQ Docs',
			defaultLocale: 'root',
			locales: {
				root: {
					label: 'English',
					lang: 'en',
				},
				'zh-cn': {
					label: '简体中文',
					lang: 'zh-CN',
				},
			},
			sidebar: [
				{
					label: 'Start Here',
					translations: {
						'zh-CN': '开始'
					},
					items: [
						{ label: 'Introduction', slug: 'intro', translations: { 'zh-CN': '介绍' } },
					],
				},
				{
					label: 'Guides',
					translations: {
						'zh-CN': '指南'
					},
					autogenerate: { directory: 'guides' },
				},
				{
					label: 'Reference',
					translations: {
						'zh-CN': '参考'
					},
					autogenerate: { directory: 'reference' },
				},
			],
		}),
	],
});
