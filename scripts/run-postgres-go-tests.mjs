import { spawnSync } from 'node:child_process';

const composeFile = 'compose.test.yaml';
const databaseUrl = 'postgres://codelinks:codelinks@localhost:55432/codelinks_auth_test?sslmode=disable';

function run(command, args, options = {}) {
	const result = spawnSync(command, args, {
		stdio: 'inherit',
		...options
	});
	if (result.status !== 0) {
		process.exitCode = result.status ?? 1;
		return false;
	}
	return true;
}

let shouldStop = false;
try {
	if (!run('docker', ['compose', '-f', composeFile, 'up', '-d', 'postgres-test'])) {
		process.exit(process.exitCode);
	}
	shouldStop = true;

	run('go', ['test', './apps/auth-service/backend/...'], {
		env: {
			...process.env,
			CODELINKS_TEST_DATABASE_URL: databaseUrl
		}
	});
} finally {
	if (shouldStop) {
		const down = spawnSync('docker', ['compose', '-f', composeFile, 'down', '-v'], {
			stdio: 'inherit'
		});
		if (process.exitCode === undefined && down.status !== 0) {
			process.exitCode = down.status ?? 1;
		}
	}
}
