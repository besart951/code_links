import { expect, test } from '@playwright/test';

test('has expected h1', async ({ page }) => {
	await page.goto('/');
	await expect(page.getByRole('heading', { name: 'CodeLinks', level: 1 })).toBeVisible();
	await expect(page.getByRole('heading', { name: 'Produktzugang' })).toBeVisible();
});

test('serves the English locale', async ({ page }) => {
	await page.goto('/en');
	await expect(page.getByRole('heading', { name: 'Product access' })).toBeVisible();
});
