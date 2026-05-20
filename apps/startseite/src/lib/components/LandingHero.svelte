<script lang="ts">
  interface HeroAction {
    label: string;
    href: string;
    external?: boolean;
  }

  interface Props {
    eyebrow: string;
    title: string;
    description: string;
    primary: HeroAction;
    secondary: HeroAction;
  }

  let { eyebrow, title, description, primary, secondary }: Props = $props();
</script>

<section
  class="relative grid min-h-[82vh] items-center overflow-hidden bg-[var(--hero-bg)] px-5 py-24 sm:px-10 lg:min-h-[84vh] lg:px-20"
>
  <div class="hero-media absolute inset-0" aria-hidden="true">
    <div class="hero-grid absolute inset-0 motion-safe:animate-[grid-drift_14s_linear_infinite]"></div>
    <div class="grid-charge absolute left-0 top-0 h-[38vh] w-[36vw] min-w-64 motion-safe:animate-[grid-charge_5.8s_ease-in-out_infinite]"></div>
    <div class="lightning absolute motion-safe:animate-[lightning-drop_5.8s_ease-in-out_infinite]"></div>
    <div class="absolute right-[8vw] top-[14vh] h-[34vw] max-h-96 w-[34vw] max-w-96 rotate-12 border border-white/35"></div>
    <div class="absolute bottom-[-12vh] right-[26vw] h-[28vw] max-h-80 w-[28vw] max-w-80 -rotate-6 border border-white/20"></div>
  </div>

  <div
    class="relative z-10 w-full max-w-6xl pt-8 text-[var(--hero-text)] motion-safe:animate-[hero-reveal_700ms_cubic-bezier(0.16,1,0.3,1)_both]"
  >
    <p class="mb-5 text-sm font-black uppercase tracking-normal sm:text-base">{eyebrow}</p>
    <h1 class="text-balance text-6xl font-black leading-[0.88] tracking-normal sm:text-8xl lg:text-[132px]">
      {title}
    </h1>
    <p class="mt-7 max-w-full break-words text-lg leading-8 sm:max-w-2xl sm:text-2xl">{description}</p>
    <div class="mt-9 flex flex-wrap gap-3">
      <a
        class="inline-flex min-h-11 items-center rounded-lg bg-white px-5 text-sm font-black text-[#14211f] shadow-2xl shadow-black/20 transition hover:-translate-y-0.5 hover:bg-[#eef7f3] focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-[#14211f]"
        href={primary.href}
        target={primary.external ? '_blank' : undefined}
        rel={primary.external ? 'noreferrer' : undefined}
      >
        {primary.label}
      </a>
      <a
        class="inline-flex min-h-11 items-center rounded-lg border border-white/45 bg-white/15 px-5 text-sm font-black text-white backdrop-blur transition hover:-translate-y-0.5 hover:bg-white/25 focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-[#14211f]"
        href={secondary.href}
        target={secondary.external ? '_blank' : undefined}
        rel={secondary.external ? 'noreferrer' : undefined}
      >
        {secondary.label}
      </a>
    </div>
  </div>
</section>

<style>
  .hero-media {
    background: var(--hero-media);
  }

  .hero-grid {
    background-image:
      linear-gradient(rgba(255, 255, 255, 0.16) 1px, transparent 1px),
      linear-gradient(90deg, rgba(255, 255, 255, 0.16) 1px, transparent 1px);
    background-size: 52px 52px;
    mask-image: linear-gradient(90deg, black, rgba(0, 0, 0, 0.82) 42%, transparent 82%);
  }

  .grid-charge {
    background:
      linear-gradient(180deg, rgba(255, 255, 255, 0), rgba(0, 213, 255, 0.24), rgba(255, 255, 255, 0)),
      linear-gradient(90deg, rgba(255, 255, 255, 0.5), rgba(0, 213, 255, 0.02));
    filter: blur(1px);
    mix-blend-mode: screen;
  }

  .lightning {
    left: clamp(18px, 6vw, 88px);
    top: -16vh;
    width: clamp(92px, 12vw, 160px);
    height: 58vh;
    opacity: 0;
    filter: drop-shadow(0 0 14px rgba(0, 213, 255, 0.95))
      drop-shadow(0 0 34px rgba(255, 255, 255, 0.58));
  }

  .lightning::before,
  .lightning::after {
    position: absolute;
    inset: 0;
    content: '';
    clip-path: polygon(
      48% 0,
      66% 0,
      53% 25%,
      76% 25%,
      47% 58%,
      63% 58%,
      31% 100%,
      42% 66%,
      22% 66%,
      45% 34%,
      29% 34%
    );
  }

  .lightning::before {
    background: linear-gradient(180deg, #ffffff 0%, var(--accent-strong) 48%, #ffffff 100%);
  }

  .lightning::after {
    scale: 1.35 1.05;
    background: rgba(0, 213, 255, 0.2);
    filter: blur(16px);
  }
</style>
