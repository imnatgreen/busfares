import { sveltekit } from '@sveltejs/kit/vite';
import { isoImport } from 'vite-plugin-iso-import'
import type { UserConfig } from 'vite';

const config: UserConfig = {
	plugins: [sveltekit(), isoImport()],
};

export default config;
