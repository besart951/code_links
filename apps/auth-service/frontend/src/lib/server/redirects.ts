export function safeRedirectTo(value: FormDataEntryValue | string | null, fallback = '/') {
	const redirectTo = typeof value === 'string' ? value : null;
	if (!redirectTo) return fallback;
	if (redirectTo.startsWith('/') && !redirectTo.startsWith('//')) return redirectTo;

	try {
		const url = new URL(redirectTo);
		if (
			url.hostname.endsWith('codelinks.localhost') ||
			url.hostname === 'localhost' ||
			url.hostname === '127.0.0.1'
		) {
			return url.toString();
		}
	} catch {
		return fallback;
	}

	return fallback;
}
