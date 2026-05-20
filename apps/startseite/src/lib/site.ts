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
  productSelectionEyebrow: string;
  productSelectionTitle: string;
  productSelectionDescription: string;
  selectedProductLabel: string;
  pricingEyebrow: string;
  pricingTitle: string;
  pricingDescription: string;
  includedLabel: string;
  priceFromLabel: string;
  pricePerMonthLabel: string;
  footerTagline: string;
  footerProductsTitle: string;
  footerPlatformTitle: string;
  footerLegalTitle: string;
  footerContact: string;
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
    productSelectionEyebrow: 'Produktauswahl',
    productSelectionTitle: 'Wähle den passenden Einstieg in die CodeLinks Suite.',
    productSelectionDescription:
      'Jedes Produkt bleibt fokussiert, nutzt aber dieselbe Platform für Login, Mandanten und Abos.',
    selectedProductLabel: 'Ausgewähltes Produkt',
    pricingEyebrow: 'Preisübersicht',
    pricingTitle: 'Startklar für einzelne Teams, erweiterbar für mehrere Produkte.',
    pricingDescription:
      'Die Preise zeigen typische Einstiegspunkte. Grössere Mandanten und Produktbündel werden über die Platform freigeschaltet.',
    includedLabel: 'Enthalten',
    priceFromLabel: 'ab',
    pricePerMonthLabel: 'pro Monat',
    footerTagline:
      'CodeLinks bündelt Fachprodukte mit gemeinsamer Identität, Mandantensteuerung und Entitlements.',
    footerProductsTitle: 'Produkte',
    footerPlatformTitle: 'Platform',
    footerLegalTitle: 'Rechtliches',
    footerContact: 'Kontakt',
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
    productSelectionEyebrow: 'Product selection',
    productSelectionTitle: 'Choose the right entry point into the CodeLinks suite.',
    productSelectionDescription:
      'Each Product stays focused while sharing the same Platform for login, Tenants and subscriptions.',
    selectedProductLabel: 'Selected Product',
    pricingEyebrow: 'Price overview',
    pricingTitle: 'Ready for individual teams, expandable across Products.',
    pricingDescription:
      'Prices show typical starting points. Larger Tenants and bundles are unlocked through the Platform.',
    includedLabel: 'Included',
    priceFromLabel: 'from',
    pricePerMonthLabel: 'per month',
    footerTagline:
      'CodeLinks combines specialist Products with shared identity, Tenant control and Entitlements.',
    footerProductsTitle: 'Products',
    footerPlatformTitle: 'Platform',
    footerLegalTitle: 'Legal',
    footerContact: 'Contact',
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
    productSelectionEyebrow: 'Sélection de produits',
    productSelectionTitle: "Choisissez le bon point d'entrée dans la suite CodeLinks.",
    productSelectionDescription:
      'Chaque produit reste ciblé tout en partageant la même Platform pour la connexion, les mandants et les abonnements.',
    selectedProductLabel: 'Produit sélectionné',
    pricingEyebrow: 'Aperçu des prix',
    pricingTitle: 'Prêt pour les équipes individuelles, extensible à plusieurs produits.',
    pricingDescription:
      'Les prix indiquent des points de départ typiques. Les grands mandants et bundles sont activés via la Platform.',
    includedLabel: 'Inclus',
    priceFromLabel: 'dès',
    pricePerMonthLabel: 'par mois',
    footerTagline:
      'CodeLinks associe des produits spécialisés à une identité commune, une gestion des mandants et des Entitlements.',
    footerProductsTitle: 'Produits',
    footerPlatformTitle: 'Platform',
    footerLegalTitle: 'Légal',
    footerContact: 'Contact',
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
    productSelectionEyebrow: 'Selezione prodotto',
    productSelectionTitle: "Scegli il punto d'ingresso giusto nella suite CodeLinks.",
    productSelectionDescription:
      'Ogni prodotto resta focalizzato, ma condivide la stessa Platform per login, tenant e abbonamenti.',
    selectedProductLabel: 'Prodotto selezionato',
    pricingEyebrow: 'Panoramica prezzi',
    pricingTitle: 'Pronto per singoli team, estendibile a più prodotti.',
    pricingDescription:
      'I prezzi mostrano punti di partenza tipici. Tenant più grandi e bundle vengono attivati tramite la Platform.',
    includedLabel: 'Incluso',
    priceFromLabel: 'da',
    pricePerMonthLabel: 'al mese',
    footerTagline:
      'CodeLinks combina prodotti specialistici con identità condivisa, gestione tenant ed Entitlements.',
    footerProductsTitle: 'Prodotti',
    footerPlatformTitle: 'Platform',
    footerLegalTitle: 'Legale',
    footerContact: 'Contatto',
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
    productSelectionEyebrow: 'Selección de productos',
    productSelectionTitle: 'Elige el punto de entrada adecuado a la suite CodeLinks.',
    productSelectionDescription:
      'Cada producto mantiene su foco y comparte la misma Platform para inicio de sesión, tenants y suscripciones.',
    selectedProductLabel: 'Producto seleccionado',
    pricingEyebrow: 'Resumen de precios',
    pricingTitle: 'Listo para equipos individuales, ampliable a varios productos.',
    pricingDescription:
      'Los precios muestran puntos de partida típicos. Tenants más grandes y bundles se activan mediante la Platform.',
    includedLabel: 'Incluido',
    priceFromLabel: 'desde',
    pricePerMonthLabel: 'al mes',
    footerTagline:
      'CodeLinks combina productos especializados con identidad compartida, gestión de tenants y Entitlements.',
    footerProductsTitle: 'Productos',
    footerPlatformTitle: 'Platform',
    footerLegalTitle: 'Legal',
    footerContact: 'Contacto',
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

export type PricePlanContent = {
  name: string;
  audience: string;
  price: string;
  summary: string;
  features: string[];
  cta: string;
};

export type PricePlan = {
  id: string;
  productSlug?: string;
  highlighted?: boolean;
  content: Record<Locale, PricePlanContent>;
};

export type LocalizedPricePlan = Omit<PricePlan, 'content'> & PricePlanContent;

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

export const pricingPlans: PricePlan[] = [
  {
    id: 'planer-starter',
    productSlug: 'planer-link',
    content: {
      de: {
        name: 'PlanerLink Starter',
        audience: 'Für kleine Teams',
        price: 'CHF 19',
        summary: 'Planung, Exporte und synchronisierte Einsatzdaten für den täglichen Ablauf.',
        features: ['1 Tenant', 'Dienstplanung', 'PDF- und Excel-Export', 'Basis-Support'],
        cta: 'PlanerLink ansehen'
      },
      en: {
        name: 'PlanerLink Starter',
        audience: 'For small teams',
        price: 'CHF 19',
        summary: 'Planning, exports and synchronized assignment data for daily operations.',
        features: ['1 Tenant', 'Duty planning', 'PDF and Excel export', 'Basic support'],
        cta: 'View PlanerLink'
      },
      fr: {
        name: 'PlanerLink Starter',
        audience: 'Pour petites équipes',
        price: 'CHF 19',
        summary:
          "Planification, exports et données d'intervention synchronisées pour le quotidien.",
        features: ['1 mandant', 'Planification', 'Export PDF et Excel', 'Support de base'],
        cta: 'Voir PlanerLink'
      },
      it: {
        name: 'PlanerLink Starter',
        audience: 'Per piccoli team',
        price: 'CHF 19',
        summary: 'Pianificazione, export e dati intervento sincronizzati per il lavoro quotidiano.',
        features: ['1 tenant', 'Pianificazione turni', 'Export PDF ed Excel', 'Supporto base'],
        cta: 'Vedi PlanerLink'
      },
      es: {
        name: 'PlanerLink Starter',
        audience: 'Para equipos pequeños',
        price: 'CHF 19',
        summary: 'Planificación, exportaciones y datos de servicio sincronizados para el día a día.',
        features: ['1 tenant', 'Planificación de turnos', 'Exportación PDF y Excel', 'Soporte básico'],
        cta: 'Ver PlanerLink'
      }
    }
  },
  {
    id: 'infra-pro',
    productSlug: 'infra-link',
    highlighted: true,
    content: {
      de: {
        name: 'InfraLink Pro',
        audience: 'Für technische Betreiber',
        price: 'CHF 49',
        summary: 'Strukturierte Gebäudeautomation mit Berechtigungen, Katalogen und Projektdaten.',
        features: ['BACnet-Katalog', 'SPS- und Feldgeräte', 'Projektrollen', 'FeatureAccess-Steuerung'],
        cta: 'InfraLink ansehen'
      },
      en: {
        name: 'InfraLink Pro',
        audience: 'For technical operators',
        price: 'CHF 49',
        summary: 'Structured building automation with permissions, catalogs and project data.',
        features: ['BACnet catalog', 'PLCs and field devices', 'Project roles', 'FeatureAccess control'],
        cta: 'View InfraLink'
      },
      fr: {
        name: 'InfraLink Pro',
        audience: 'Pour exploitants techniques',
        price: 'CHF 49',
        summary:
          'Automatisation du bâtiment structurée avec droits, catalogues et données de projet.',
        features: ['Catalogue BACnet', 'Automates et terrain', 'Rôles projet', 'Gestion FeatureAccess'],
        cta: 'Voir InfraLink'
      },
      it: {
        name: 'InfraLink Pro',
        audience: 'Per gestori tecnici',
        price: 'CHF 49',
        summary: 'Automazione edifici strutturata con permessi, cataloghi e dati progetto.',
        features: ['Catalogo BACnet', 'PLC e campo', 'Ruoli progetto', 'Controllo FeatureAccess'],
        cta: 'Vedi InfraLink'
      },
      es: {
        name: 'InfraLink Pro',
        audience: 'Para operadores técnicos',
        price: 'CHF 49',
        summary: 'Automatización de edificios estructurada con permisos, catálogos y datos de proyecto.',
        features: ['Catálogo BACnet', 'PLC y campo', 'Roles de proyecto', 'Control FeatureAccess'],
        cta: 'Ver InfraLink'
      }
    }
  },
  {
    id: 'platform-bundle',
    content: {
      de: {
        name: 'CodeLinks Bundle',
        audience: 'Für mehrere Produkte',
        price: 'Individuell',
        summary: 'Gemeinsame Platform, Produktzugriff, Entitlements und Planlogik für Mandanten.',
        features: ['Mehrere Products', 'ProductAccess', 'Entitlements', 'Individuelle FeatureLimits'],
        cta: 'Kontakt aufnehmen'
      },
      en: {
        name: 'CodeLinks Bundle',
        audience: 'For multiple Products',
        price: 'Custom',
        summary: 'Shared Platform, ProductAccess, Entitlements and Plan logic for Tenants.',
        features: ['Multiple Products', 'ProductAccess', 'Entitlements', 'Custom FeatureLimits'],
        cta: 'Contact us'
      },
      fr: {
        name: 'Bundle CodeLinks',
        audience: 'Pour plusieurs produits',
        price: 'Sur mesure',
        summary:
          'Platform commune, ProductAccess, Entitlements et logique de Plans pour les mandants.',
        features: ['Plusieurs produits', 'ProductAccess', 'Entitlements', 'FeatureLimits sur mesure'],
        cta: 'Nous contacter'
      },
      it: {
        name: 'Bundle CodeLinks',
        audience: 'Per più prodotti',
        price: 'Su misura',
        summary: 'Platform condivisa, ProductAccess, Entitlements e logica Plan per i tenant.',
        features: ['Più prodotti', 'ProductAccess', 'Entitlements', 'FeatureLimits su misura'],
        cta: 'Contattaci'
      },
      es: {
        name: 'Bundle CodeLinks',
        audience: 'Para varios productos',
        price: 'Personalizado',
        summary: 'Platform compartida, ProductAccess, Entitlements y lógica de Plan para tenants.',
        features: ['Varios productos', 'ProductAccess', 'Entitlements', 'FeatureLimits personalizados'],
        cta: 'Contactar'
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

export function getPricingPlans(locale: Locale): LocalizedPricePlan[] {
  return pricingPlans.map((plan) => localizePricePlan(plan, locale));
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

function localizeProduct(product: ProductPage, locale: Locale): LocalizedProductPage {
  const { content, ...base } = product;
  return {
    ...base,
    ...content[locale]
  };
}

function localizePricePlan(plan: PricePlan, locale: Locale): LocalizedPricePlan {
  const { content, ...base } = plan;
  return {
    ...base,
    ...content[locale]
  };
}
