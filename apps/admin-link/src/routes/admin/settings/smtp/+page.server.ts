import { fail } from '@sveltejs/kit';
import { hasPermission } from '$lib/domain/admin-access/permissions';
import type { SmtpEncryption } from '$lib/domain/smtp/types';
import { formString } from '$lib/server/admin-route-helpers';
import { createAdminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

function parseEncryption(value: string): SmtpEncryption {
	return value === 'none' || value === 'ssl' || value === 'tls' || value === 'starttls' ? value : 'starttls';
}

export async function load(event) {
	const admin = requireAdmin(event.locals);
	const adminContainer = createAdminContainer(event);

	return {
		settings: await adminContainer.getSmtpSettings.execute(admin),
		canUpdate: hasPermission(admin, 'admin.smtp_settings.update')
	};
}

export const actions = {
	save: async (event) => {
		const admin = requireAdmin(event.locals);
		const adminContainer = createAdminContainer(event);
		const formData = await event.request.formData();
		const host = formString(formData, 'host');
		const port = Number(formData.get('port') ?? 0);
		const fromEmail = formString(formData, 'fromEmail');
		const replyToEmail = formString(formData, 'replyToEmail');

		if (!host || !port || !fromEmail || !replyToEmail) {
			return fail(400, { error: true, message: 'Pflichtfelder fehlen.' });
		}

		await adminContainer.updateSmtpSettings.execute(admin, {
			host,
			port,
			username: formString(formData, 'username'),
			password: String(formData.get('password') ?? ''),
			encryption: parseEncryption(String(formData.get('encryption') ?? '')),
			fromEmail,
			fromName: formString(formData, 'fromName'),
			replyToEmail,
			active: formData.get('active') === 'on'
		});

		return { ok: true, message: 'SMTP-Einstellungen gespeichert.' };
	},
	testEmail: async (event) => {
		const admin = requireAdmin(event.locals);
		const adminContainer = createAdminContainer(event);
		const formData = await event.request.formData();
		const recipient = formString(formData, 'recipient');

		if (!recipient) {
			return fail(400, { error: true, message: 'Empfänger fehlt.' });
		}

		await adminContainer.sendTestEmail.execute(admin, recipient);
		return { ok: true, message: 'Test-E-Mail gesendet.' };
	}
};
