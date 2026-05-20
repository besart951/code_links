import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, '.', '');
	const platformUrl = env.PLATFORM_API_URL ?? 'http://localhost:8081';

	return {
		plugins: [tailwindcss(), sveltekit()],
		test: {
			environment: 'jsdom',
			include: ['src/**/*.test.ts']
		},
		server: {
			proxy: {
				'/api': {
					target: platformUrl,
					changeOrigin: true,
					secure: false
				}
			}
		}
	};
});
