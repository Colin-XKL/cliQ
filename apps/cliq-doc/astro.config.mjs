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
					label: '简体中文',
					lang: 'zh-CN',
				},
				en: {
					label: 'English',
					lang: 'en',
				},
			},
			sidebar: [
				{
					label: '开始',
					translations: {
						'en': 'Start Here'
					},
					items: [
						{ label: '介绍', slug: 'intro', translations: { 'en': 'Introduction' } },
					],
				},
				{
					label: '指南',
					translations: {
						'en': 'Guides'
					},
					autogenerate: { directory: 'guides' },
				},
				{
					label: '参考',
					translations: {
						'en': 'Reference'
					},
					autogenerate: { directory: 'reference' },
				},
			],
		}),
	],
});
