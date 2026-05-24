import { spawn } from 'node:child_process';
import path from 'node:path';
import process from 'node:process';

const isWindows = process.platform === 'win32';
const args = [
	'--filter',
	'@codelinks/admin-link',
	'exec',
	'vite',
	'dev',
	'--host',
	'localhost',
	'--port',
	'5178'
];
const env = {
	...process.env,
	ADMIN_LINK_MOCK_AUTH: process.env.ADMIN_LINK_MOCK_AUTH ?? 'true',
	ADMIN_LINK_DATA_SOURCE: process.env.ADMIN_LINK_DATA_SOURCE ?? 'mock'
};

const pnpmExecPath = process.env.npm_execpath;
const canReusePnpm = pnpmExecPath?.toLowerCase().includes('pnpm') ?? false;
const pnpmExecExtension = canReusePnpm ? path.extname(pnpmExecPath).toLowerCase() : '';
const pnpmExecIsJS = ['.js', '.cjs', '.mjs'].includes(pnpmExecExtension);
const pnpmExecNeedsShell = ['.cmd', '.bat'].includes(pnpmExecExtension);

const child = canReusePnpm
	? pnpmExecIsJS
		? spawn(process.execPath, [pnpmExecPath, ...args], { stdio: 'inherit', env })
		: spawn(pnpmExecPath, args, { stdio: 'inherit', env, shell: pnpmExecNeedsShell })
	: isWindows
		? spawn('pnpm', args, { stdio: 'inherit', shell: true, env })
		: spawn('pnpm', args, { stdio: 'inherit', env });

for (const signal of ['SIGINT', 'SIGTERM']) {
	process.on(signal, () => {
		child.kill(signal);
	});
}

child.on('error', (error) => {
	console.error(error);
	process.exit(1);
});

child.on('exit', (code, signal) => {
	if (signal) {
		process.exit(1);
	}

	process.exit(code ?? 0);
});
