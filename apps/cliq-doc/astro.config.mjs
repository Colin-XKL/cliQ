// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'cliQ Docs',
			defaultLocale: 'zh-CN',
			locales: {
				'zh-CN': {
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
					label: 'Start Here',
					translations: {
						'zh-CN': '开始'
					},
					autogenerate: {
						directory: 'start'
					}
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
