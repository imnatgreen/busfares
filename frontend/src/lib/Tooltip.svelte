<!-- https://github.com/dimfeld/svelte-leaflet-demo/blob/master/full/src/map/Popup.svelte -->
<script lang="ts">
  import L from 'leaflet?client';
  import { getContext } from 'svelte';

  let classNames: string | undefined = undefined;
  export { classNames as class };

  export let permanent = false;
  export let sticky = false;
  export let interactive = false;
  export let direction: L.Direction = 'auto';
  export let opacity = 0.9;

  export let tooltip: L.Tooltip | undefined = undefined;

  let showContents = permanent;
  let tooltipOpen = permanent;

  const layer = getContext<() => L.Layer>('layer')();
  function createPopup(tooltipElement: HTMLElement) {
    tooltip = L.tooltip({permanent, sticky, interactive, direction, opacity}).setContent(tooltipElement);
    layer.bindTooltip(tooltip);

    layer.on('tooltipopen', () => {
      tooltipOpen = true;
      showContents = true;
    });

    layer.on('tooltipclose', () => {
      tooltipOpen = false;
      // Wait for the tooltip to completely fade out before destroying it.
      // Otherwise the fade out looks weird as the contents disappear too early.
      setTimeout(() => {
        if (!tooltipOpen) {
          showContents = false;
        }
      }, 500);
    });

    return {
      destroy() {
        if (tooltip) {
          layer.unbindTooltip();
          tooltip.remove();
          tooltip = undefined;
        }
      },
    };
  }
</script>

<div class="hidden">
  <div use:createPopup class={classNames}>
    {#if showContents}
      <slot />
    {/if}
  </div>
</div>