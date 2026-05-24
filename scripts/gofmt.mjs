import { execFileSync } from 'node:child_process';
import { existsSync } from 'node:fs';
import { resolve } from 'node:path';

const root = resolve(import.meta.dirname, '..');
const check = process.argv.includes('--check');

const targets = [
	'apps/auth-service/backend',
	'apps/infra-link/backend',
	'apps/planer-link/backend',
	'apps/loka-link/backend',
	'packages/adminaccess',
	'packages/productauth',
	'packages/productcatalog',
	'packages/productserver'
].filter((path) => existsSync(resolve(root, path)));

const args = check ? ['-l', ...targets] : ['-w', ...targets];
const output = execFileSync('gofmt', args, { cwd: root, encoding: 'utf8' });

if (check && output.trim()) {
	process.stdout.write(output);
	process.exitCode = 1;
}
