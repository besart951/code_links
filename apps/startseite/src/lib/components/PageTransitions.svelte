<script lang="ts">
  import { browser } from '$app/environment';
  import { onNavigate } from '$app/navigation';

  type ViewTransitionDocument = Document & {
    startViewTransition?: (callback: () => Promise<void> | void) => { finished: Promise<void> };
  };

  onNavigate((navigation) => {
    if (!browser) {
      return;
    }

    const transitionDocument = document as ViewTransitionDocument;

    if (!transitionDocument.startViewTransition) {
      return;
    }

    return new Promise<void>((resolve) => {
      transitionDocument.startViewTransition(async () => {
        resolve();
        await navigation.complete;
      });
    });
  });
</script>
