<script lang="ts">
  const products = [
    { name: 'infra_link', href: 'https://infra.codelinks.ch', text: 'Gebaeudeautomation, BACnet, SPS und Feldgeraete.' },
    { name: 'planer_link', href: 'https://planer.codelinks.ch', text: 'Dienstplanung, Tourenplanung, PDF und Excel.' },
    { name: 'loka_link', href: 'https://loka.codelinks.ch', text: 'Naechstes Produkt im CodeLinks Verbund.' }
  ];
</script>

<svelte:head>
  <title>CodeLinks</title>
  <meta
    name="description"
    content="CodeLinks ist die gemeinsame Plattform fuer spezialisierte Fachprodukte."
  />
</svelte:head>

<main class="page">
  <section class="hero">
    <div class="hero__media" aria-hidden="true">
      <div class="signal signal--one"></div>
      <div class="signal signal--two"></div>
      <div class="grid"></div>
    </div>
    <div class="hero__content">
      <p class="eyebrow">CodeLinks Plattform</p>
      <h1>CodeLinks</h1>
      <p class="lead">
        Ein Monorepo fuer eigenstaendige Produkte mit gemeinsamer Identitaet,
        Mandantensteuerung und Abo-basierten Features.
      </p>
      <div class="actions">
        <a href="https://auth.codelinks.ch">Anmelden</a>
        <a href="https://infra.codelinks.ch" class="secondary">Produkte oeffnen</a>
      </div>
    </div>
  </section>

  <section class="products" aria-label="Produkte">
    {#each products as product}
      <a class="product" href={product.href}>
        <span>{product.name}</span>
        <p>{product.text}</p>
      </a>
    {/each}
  </section>
</main>

<style>
  :global(body) {
    margin: 0;
    min-width: 320px;
    font-family:
      Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
      sans-serif;
    color: #14211f;
    background: #f4f7f5;
  }

  .page {
    min-height: 100vh;
  }

  .hero {
    position: relative;
    min-height: 78vh;
    overflow: hidden;
    display: grid;
    place-items: center start;
    padding: clamp(32px, 8vw, 96px);
    background: #dfe9e4;
  }

  .hero__media {
    position: absolute;
    inset: 0;
    background:
      linear-gradient(120deg, rgba(20, 33, 31, 0.9), rgba(20, 33, 31, 0.18)),
      radial-gradient(circle at 75% 25%, #78b7a2, transparent 28%),
      linear-gradient(135deg, #173b36, #e2ede7);
  }

  .grid {
    position: absolute;
    inset: 0;
    background-image:
      linear-gradient(rgba(255, 255, 255, 0.14) 1px, transparent 1px),
      linear-gradient(90deg, rgba(255, 255, 255, 0.14) 1px, transparent 1px);
    background-size: 56px 56px;
    mask-image: linear-gradient(90deg, black, transparent 75%);
  }

  .signal {
    position: absolute;
    width: 34vw;
    aspect-ratio: 1;
    border: 1px solid rgba(255, 255, 255, 0.36);
    transform: rotate(18deg);
  }

  .signal--one {
    right: 10vw;
    top: 12vh;
  }

  .signal--two {
    right: 24vw;
    bottom: -8vh;
  }

  .hero__content {
    position: relative;
    max-width: 760px;
    color: white;
  }

  .eyebrow {
    margin: 0 0 18px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0;
  }

  h1 {
    margin: 0;
    font-size: clamp(56px, 12vw, 132px);
    line-height: 0.88;
    letter-spacing: 0;
  }

  .lead {
    max-width: 620px;
    margin: 28px 0 0;
    font-size: clamp(18px, 2vw, 24px);
    line-height: 1.45;
  }

  .actions {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    margin-top: 34px;
  }

  .actions a {
    display: inline-flex;
    min-height: 44px;
    align-items: center;
    border-radius: 8px;
    padding: 0 18px;
    background: #ffffff;
    color: #14211f;
    font-weight: 700;
    text-decoration: none;
  }

  .actions .secondary {
    background: rgba(255, 255, 255, 0.16);
    color: #ffffff;
    border: 1px solid rgba(255, 255, 255, 0.4);
  }

  .products {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 16px;
    max-width: 1120px;
    margin: -56px auto 0;
    padding: 0 24px 64px;
    position: relative;
  }

  .product {
    min-height: 132px;
    border-radius: 8px;
    border: 1px solid #d7dfda;
    background: #ffffff;
    padding: 22px;
    color: #14211f;
    text-decoration: none;
    box-shadow: 0 18px 40px rgba(21, 45, 40, 0.12);
  }

  .product span {
    display: block;
    font-size: 20px;
    font-weight: 800;
  }

  .product p {
    margin: 12px 0 0;
    color: #52615d;
    line-height: 1.5;
  }

  @media (max-width: 680px) {
    .hero {
      align-items: end;
      min-height: 72vh;
      padding: 28px 22px 86px;
    }

    .products {
      margin-top: -40px;
    }
  }
</style>
