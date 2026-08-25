import { expect, test } from '@playwright/test';

test('renders the German landing page', async ({ page }) => {
	await page.goto('/');

	await expect(page.getByRole('heading', { name: 'CodeLinks', level: 1 })).toBeVisible();
	await expect(page.getByRole('heading', { name: 'Wer wir sind' })).toBeVisible();
	await expect(page.getByRole('heading', { name: 'Produktzugang' })).toBeVisible();
});

test('serves the English locale', async ({ page }) => {
	await page.goto('/en');

	await expect(page.getByRole('heading', { name: 'Who we are' })).toBeVisible();
	await expect(page.getByRole('heading', { name: 'Product access' })).toBeVisible();
});

test('language selector routes to localized pages', async ({ page }) => {
	await page.goto('/');

	await page.getByRole('button', { name: 'Sprache' }).click();
	await page.getByRole('link', { name: 'English' }).click();

	await expect(page).toHaveURL(/\/en\/?$/);
	await expect(page.getByRole('heading', { name: 'Who we are' })).toBeVisible();
});

test('theme selector applies light and dark modes', async ({ page }) => {
	await page.goto('/');

	await page.getByRole('button', { name: 'Darstellung' }).click();
	await page.getByText('Dunkel').click();
	await expect(page.locator('html')).toHaveClass(/dark/);

	await page.getByRole('button', { name: 'Darstellung' }).click();
	await page.getByText('Hell').click();
	await expect(page.locator('html')).not.toHaveClass(/dark/);
});

test('login link points to the auth app with a return URL', async ({ page }) => {
	await page.goto('/');

	await expect(page.getByRole('navigation').getByRole('link', { name: 'Login', exact: true })).toHaveAttribute(
		'href',
		/http:\/\/localhost:5174\/login\?redirectTo=http%3A%2F%2Flocalhost%3A4173%2F/
	);
});

test('mobile navigation opens and keeps landing links reachable', async ({ page }) => {
	await page.setViewportSize({ width: 390, height: 844 });
	await page.goto('/');

	await page.getByRole('button', { name: 'Navigation oeffnen' }).click();

	await expect(page.getByRole('link', { name: 'Motivation' })).toBeVisible();
	await expect(page.getByRole('link', { name: 'Produkte', exact: true })).toBeVisible();
});
