export const locales = ['de', 'en', 'fr', 'it', 'es'] as const;
export type Locale = (typeof locales)[number];

export const defaultLocale: Locale = 'de';

export const site = {
  title: 'CodeLinks',
  url: 'https://codelinks.ch'
} as const;

export const localeLabels: Record<Locale, string> = {
  de: 'DE',
  en: 'EN',
  fr: 'FR',
  it: 'IT',
  es: 'ES'
};

export const localeNames: Record<Locale, string> = {
  de: 'Deutsch',
  en: 'English',
  fr: 'Français',
  it: 'Italiano',
  es: 'Español'
};

export const htmlLang: Record<Locale, string> = {
  de: 'de-CH',
  en: 'en',
  fr: 'fr',
  it: 'it',
  es: 'es'
};

export type SiteCopy = {
  title: string;
  description: string;
  eyebrow: string;
  signIn: string;
  viewProducts: string;
  productsLabel: string;
  homeLabel: string;
  openProduct: string;
  plannedAction: string;
  featuresLabel: string;
  controlsLabel: string;
  languageLabel: string;
  themeLabel: string;
  lightTheme: string;
  darkTheme: string;
};

export const copy: Record<Locale, SiteCopy> = {
  de: {
    title: 'CodeLinks',
    description:
      'CodeLinks verbindet eigenständige Fachprodukte mit gemeinsamer Identität, Mandantensteuerung und Abo-basierten Features.',
    eyebrow: 'CodeLinks Plattform',
    signIn: 'Anmelden',
    viewProducts: 'Produkte ansehen',
    productsLabel: 'Produkte',
    homeLabel: 'CodeLinks Startseite',
    openProduct: 'Produkt öffnen',
    plannedAction: 'Vormerken',
    featuresLabel: 'Funktionen',
    controlsLabel: 'Sprache und Darstellung',
    languageLabel: 'Sprache',
    themeLabel: 'Darstellung',
    lightTheme: 'Hell',
    darkTheme: 'Dunkel'
  },
  en: {
    title: 'CodeLinks',
    description:
      'CodeLinks connects independent specialist products with shared identity, tenant control and subscription-based features.',
    eyebrow: 'CodeLinks Platform',
    signIn: 'Sign in',
    viewProducts: 'View products',
    productsLabel: 'Products',
    homeLabel: 'CodeLinks home',
    openProduct: 'Open product',
    plannedAction: 'Follow launch',
    featuresLabel: 'Features',
    controlsLabel: 'Language and appearance',
    languageLabel: 'Language',
    themeLabel: 'Appearance',
    lightTheme: 'Light',
    darkTheme: 'Dark'
  },
  fr: {
    title: 'CodeLinks',
    description:
      'CodeLinks relie des produits spécialisés indépendants avec une identité commune, une gestion des mandants et des fonctionnalités par abonnement.',
    eyebrow: 'Plateforme CodeLinks',
    signIn: 'Se connecter',
    viewProducts: 'Voir les produits',
    productsLabel: 'Produits',
    homeLabel: 'Accueil CodeLinks',
    openProduct: 'Ouvrir le produit',
    plannedAction: 'Suivre le lancement',
    featuresLabel: 'Fonctionnalités',
    controlsLabel: 'Langue et apparence',
    languageLabel: 'Langue',
    themeLabel: 'Apparence',
    lightTheme: 'Clair',
    darkTheme: 'Sombre'
  },
  it: {
    title: 'CodeLinks',
    description:
      'CodeLinks collega prodotti specialistici indipendenti con identità condivisa, gestione dei tenant e funzionalità basate su abbonamento.',
    eyebrow: 'Piattaforma CodeLinks',
    signIn: 'Accedi',
    viewProducts: 'Vedi prodotti',
    productsLabel: 'Prodotti',
    homeLabel: 'Home CodeLinks',
    openProduct: 'Apri prodotto',
    plannedAction: 'Segui il lancio',
    featuresLabel: 'Funzionalità',
    controlsLabel: 'Lingua e aspetto',
    languageLabel: 'Lingua',
    themeLabel: 'Aspetto',
    lightTheme: 'Chiaro',
    darkTheme: 'Scuro'
  },
  es: {
    title: 'CodeLinks',
    description:
      'CodeLinks conecta productos especializados independientes con identidad compartida, gestión de tenants y funciones basadas en suscripción.',
    eyebrow: 'Plataforma CodeLinks',
    signIn: 'Iniciar sesión',
    viewProducts: 'Ver productos',
    productsLabel: 'Productos',
    homeLabel: 'Inicio de CodeLinks',
    openProduct: 'Abrir producto',
    plannedAction: 'Seguir lanzamiento',
    featuresLabel: 'Funciones',
    controlsLabel: 'Idioma y apariencia',
    languageLabel: 'Idioma',
    themeLabel: 'Apariencia',
    lightTheme: 'Claro',
    darkTheme: 'Oscuro'
  }
};

export type ProductContent = {
  name: string;
  headline: string;
  summary: string;
  features: string[];
};

export type ProductPage = {
  slug: string;
  key: 'infra_link' | 'planer_link' | 'loka_link';
  appUrl: string;
  status: 'available' | 'planned';
  content: Record<Locale, ProductContent>;
};

export type LocalizedProductPage = Omit<ProductPage, 'content'> & ProductContent;

export const products: ProductPage[] = [
  {
    slug: 'infra-link',
    key: 'infra_link',
    appUrl: 'https://infra.codelinks.ch',
    status: 'available',
    content: {
      de: {
        name: 'InfraLink',
        headline: 'Gebäudeautomation strukturiert verwalten',
        summary:
          'InfraLink bündelt BACnet, SPS, Feldgeräte und Projektstruktur für technische Gebäudeinfrastruktur.',
        features: ['BACnet-Objekte', 'SPS-Struktur', 'Feldgeräte', 'Projektberechtigungen']
      },
      en: {
        name: 'InfraLink',
        headline: 'Structured building automation management',
        summary:
          'InfraLink brings BACnet, PLCs, field devices and project structure together for technical building infrastructure.',
        features: ['BACnet objects', 'PLC structure', 'Field devices', 'Project permissions']
      },
      fr: {
        name: 'InfraLink',
        headline: "Gérer l'automatisation du bâtiment de manière structurée",
        summary:
          "InfraLink regroupe BACnet, automates, appareils de terrain et structure de projet pour l'infrastructure technique du bâtiment.",
        features: ['Objets BACnet', 'Structure automate', 'Appareils de terrain', 'Droits de projet']
      },
      it: {
        name: 'InfraLink',
        headline: "Gestire l'automazione degli edifici in modo strutturato",
        summary:
          "InfraLink riunisce BACnet, PLC, dispositivi di campo e struttura di progetto per l'infrastruttura tecnica degli edifici.",
        features: ['Oggetti BACnet', 'Struttura PLC', 'Dispositivi di campo', 'Permessi di progetto']
      },
      es: {
        name: 'InfraLink',
        headline: 'Gestionar la automatización de edificios de forma estructurada',
        summary:
          'InfraLink reúne BACnet, PLC, dispositivos de campo y estructura de proyecto para infraestructura técnica de edificios.',
        features: ['Objetos BACnet', 'Estructura PLC', 'Dispositivos de campo', 'Permisos de proyecto']
      }
    }
  },
  {
    slug: 'planer-link',
    key: 'planer_link',
    appUrl: 'https://planer.codelinks.ch',
    status: 'available',
    content: {
      de: {
        name: 'PlanerLink',
        headline: 'Dienst- und Einsatzplanung für kleine Teams',
        summary:
          'PlanerLink unterstützt Planung, PDF-Export, Excel-Export und synchronisierte Einsatzdaten.',
        features: ['Dienstplanung', 'PDF-Export', 'Excel-Export', 'Sync']
      },
      en: {
        name: 'PlanerLink',
        headline: 'Duty and assignment planning for small teams',
        summary:
          'PlanerLink supports planning, PDF export, Excel export and synchronized assignment data.',
        features: ['Duty planning', 'PDF export', 'Excel export', 'Sync']
      },
      fr: {
        name: 'PlanerLink',
        headline: 'Planification des services et interventions pour petites équipes',
        summary:
          "PlanerLink prend en charge la planification, l'export PDF, l'export Excel et les données d'intervention synchronisées.",
        features: ['Planification', 'Export PDF', 'Export Excel', 'Sync']
      },
      it: {
        name: 'PlanerLink',
        headline: 'Pianificazione di turni e interventi per piccoli team',
        summary:
          'PlanerLink supporta pianificazione, esportazione PDF, esportazione Excel e dati di intervento sincronizzati.',
        features: ['Pianificazione turni', 'Export PDF', 'Export Excel', 'Sync']
      },
      es: {
        name: 'PlanerLink',
        headline: 'Planificación de turnos y servicios para equipos pequeños',
        summary:
          'PlanerLink permite planificación, exportación PDF, exportación Excel y datos de servicio sincronizados.',
        features: ['Planificación de turnos', 'Exportación PDF', 'Exportación Excel', 'Sync']
      }
    }
  },
  {
    slug: 'loka-link',
    key: 'loka_link',
    appUrl: 'https://loka.codelinks.ch',
    status: 'planned',
    content: {
      de: {
        name: 'LokaLink',
        headline: 'Das nächste Produkt im CodeLinks Verbund',
        summary:
          'LokaLink ist für einen späteren Produktstart vorbereitet und nutzt dieselbe Platform-Basis.',
        features: ['Gemeinsames Login', 'Mandantenmodell', 'Abo-Steuerung']
      },
      en: {
        name: 'LokaLink',
        headline: 'The next product in the CodeLinks family',
        summary:
          'LokaLink is prepared for a later product launch and uses the same Platform foundation.',
        features: ['Shared login', 'Tenant model', 'Subscription control']
      },
      fr: {
        name: 'LokaLink',
        headline: 'Le prochain produit de la famille CodeLinks',
        summary:
          'LokaLink est préparé pour un lancement ultérieur et utilise la même base Platform.',
        features: ['Connexion commune', 'Modèle mandant', 'Gestion des abonnements']
      },
      it: {
        name: 'LokaLink',
        headline: 'Il prossimo prodotto della famiglia CodeLinks',
        summary:
          'LokaLink è preparato per un lancio successivo e usa la stessa base Platform.',
        features: ['Login condiviso', 'Modello tenant', 'Gestione abbonamenti']
      },
      es: {
        name: 'LokaLink',
        headline: 'El próximo producto de la familia CodeLinks',
        summary:
          'LokaLink está preparado para un lanzamiento posterior y usa la misma base Platform.',
        features: ['Inicio de sesión compartido', 'Modelo tenant', 'Control de suscripciones']
      }
    }
  }
];

export function parseLocale(value: string | null | undefined): Locale {
  if (value && locales.includes(value as Locale)) {
    return value as Locale;
  }

  return defaultLocale;
}

export function isLocale(value: string | null | undefined): value is Locale {
  return Boolean(value && locales.includes(value as Locale));
}

export function getSiteCopy(locale: Locale): SiteCopy {
  return copy[locale];
}

export function getProducts(locale: Locale): LocalizedProductPage[] {
  return products.map((product) => localizeProduct(product, locale));
}

export function getProductBySlug(slug: string, locale: Locale): LocalizedProductPage | undefined {
  const product = products.find((item) => item.slug === slug);
  if (!product) {
    return undefined;
  }

  return localizeProduct(product, locale);
}

export function localizedHref(path: string, locale: Locale): string {
  const normalizedPath = path === '/' ? '' : path;
  return `/${locale}${normalizedPath}`;
}

export function canonicalPathFromLocalizedPath(pathname: string): string {
  const parts = pathname.split('/').filter(Boolean);
  if (isLocale(parts[0])) {
    const rest = parts.slice(1).join('/');
    return rest ? `/${rest}` : '/';
  }

  return pathname || '/';
}

export function localeFromPath(pathname: string): Locale {
  const firstSegment = pathname.split('/').filter(Boolean)[0];
  return parseLocale(firstSegment);
}

export function absoluteUrl(path: string): string {
  return new URL(path, site.url).toString();
}

export function localizedAbsoluteUrl(path: string, locale: Locale): string {
  return absoluteUrl(localizedHref(path, locale));
}

export function alternateUrls(path: string): Array<{ locale: Locale; href: string; hrefLang: string }> {
  return locales.map((locale) => ({
    locale,
    href: localizedAbsoluteUrl(path, locale),
    hrefLang: htmlLang[locale]
  }));
}

export function canonicalContentPaths(): string[] {
  return ['/', ...products.map((product) => `/produkte/${product.slug}`)];
}

export function jsonLd(value: unknown): string {
  const json = JSON.stringify(value).replace(/</g, '\\u003c');
  return `<script type="application/ld+json">${json}</script>`;
}

function localizeProduct(product: ProductPage, locale: Locale): LocalizedProductPage {
  const { content, ...base } = product;
  return {
    ...base,
    ...content[locale]
  };
}
