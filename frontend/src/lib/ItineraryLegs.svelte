<script lang="ts">
  import {secondsToDurationString, getContrastYIQ} from '$lib/utils';
  import ItineraryLeg from '$lib/ItineraryLeg.svelte';

  export let itinerary;

  let totalWidth = itinerary.duration - itinerary.waitingTime;
  if (itinerary.legs[0].duration <= 60) {
    totalWidth -= 60;
  }

  let widths = itinerary.legs.map((leg) => {
    if (leg.duration <= 90) {
      return 0;
    } else {
      return (leg.duration / totalWidth) * 100;
    }
  });
</script>

<div class="flex flex-row flex-none gap-1 mt-1">
  {#each itinerary.legs as leg, i}
    {#if widths[i] > 0}
      <ItineraryLeg leg={leg} width={widths[i]} />
    {/if}
  {/each}
</div>