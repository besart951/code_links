import { fail } from '@sveltejs/kit';
import { hasPermission } from '$lib/domain/admin-access/permissions';
import type { SmtpEncryption } from '$lib/domain/smtp/types';
import { adminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

function parseEncryption(value: string): SmtpEncryption {
	return value === 'none' || value === 'ssl' || value === 'tls' || value === 'starttls' ? value : 'starttls';
}

export async function load({ locals }) {
	const admin = requireAdmin(locals);

	return {
		settings: await adminContainer.getSmtpSettings.execute(admin),
		canUpdate: hasPermission(admin, 'admin.smtp_settings.update')
	};
}

export const actions = {
	save: async ({ request, locals }) => {
		const admin = requireAdmin(locals);
		const formData = await request.formData();
		const host = String(formData.get('host') ?? '').trim();
		const port = Number(formData.get('port') ?? 0);
		const fromEmail = String(formData.get('fromEmail') ?? '').trim();
		const replyToEmail = String(formData.get('replyToEmail') ?? '').trim();

		if (!host || !port || !fromEmail || !replyToEmail) {
			return fail(400, { error: true, message: 'Pflichtfelder fehlen.' });
		}

		await adminContainer.updateSmtpSettings.execute(admin, {
			host,
			port,
			username: String(formData.get('username') ?? '').trim(),
			password: String(formData.get('password') ?? ''),
			encryption: parseEncryption(String(formData.get('encryption') ?? '')),
			fromEmail,
			fromName: String(formData.get('fromName') ?? '').trim(),
			replyToEmail,
			active: formData.get('active') === 'on'
		});

		return { ok: true, message: 'SMTP-Einstellungen gespeichert.' };
	},
	testEmail: async ({ request, locals }) => {
		const admin = requireAdmin(locals);
		const formData = await request.formData();
		const recipient = String(formData.get('recipient') ?? '').trim();

		if (!recipient) {
			return fail(400, { error: true, message: 'Empfänger fehlt.' });
		}

		await adminContainer.sendTestEmail.execute(admin, recipient);
		return { ok: true, message: 'Test-E-Mail gesendet.' };
	}
};
