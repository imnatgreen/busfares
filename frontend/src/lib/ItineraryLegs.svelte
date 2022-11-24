<script lang="ts">
  import {secondsToDurationString, getContrastYIQ} from '$lib/utils';

  export let itinerary;

  let widths = itinerary.legs.map((leg) => {
    return (leg.duration / (itinerary.duration - itinerary.waitingTime)) * 100;
  });
</script>

<div class="flex flex-row flex-none gap-1 mt-1">
  {#each itinerary.legs as leg, i}
    <div class="rounded-md text-sm h-8 py-1 px-2 border-gray-300 border shadow-sm transition-colors duration-300
    overflow-hidden whitespace-nowrap flex flex-row items-center flex-wrap gap-1
    {leg.mode === 'WALK' ? 'bg-gray-100' : 'bg-indigo-300'}
    {leg.colour?(getContrastYIQ(leg.colour)=='black'?'text-black':'text-white'):''}"
    style="width: {widths[i]}%;{leg.colour ? 'background-color:'+leg.colour : ''}">
      <div>
        {#if leg.mode == "BUS"}
          <svg class="w-4 h-4" fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path d="M256 0C390.4 0 480 35.2 480 80V96l0 32c17.7 0 32 14.3 32 32v64c0 17.7-14.3 32-32 32l0 160c0 17.7-14.3 32-32 32v32c0 17.7-14.3 32-32 32H384c-17.7 0-32-14.3-32-32V448H160v32c0 17.7-14.3 32-32 32H96c-17.7 0-32-14.3-32-32l0-32c-17.7 0-32-14.3-32-32l0-160c-17.7 0-32-14.3-32-32V160c0-17.7 14.3-32 32-32h0V96h0V80C32 35.2 121.6 0 256 0zM96 160v96c0 17.7 14.3 32 32 32H240V128H128c-17.7 0-32 14.3-32 32zM272 288H384c17.7 0 32-14.3 32-32V160c0-17.7-14.3-32-32-32H272V288zM112 400c17.7 0 32-14.3 32-32s-14.3-32-32-32s-32 14.3-32 32s14.3 32 32 32zm288 0c17.7 0 32-14.3 32-32s-14.3-32-32-32s-32 14.3-32 32s14.3 32 32 32zM352 80c0-8.8-7.2-16-16-16H176c-8.8 0-16 7.2-16 16s7.2 16 16 16H336c8.8 0 16-7.2 16-16z"/></svg>
        {:else if leg.mode == "WALK"}
          <svg class="w-3 h-4" fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 512"><path d="M256 48c0 26.5-21.5 48-48 48s-48-21.5-48-48s21.5-48 48-48s48 21.5 48 48zM126.5 199.3c-1 .4-1.9 .8-2.9 1.2l-8 3.5c-16.4 7.3-29 21.2-34.7 38.2l-2.6 7.8c-5.6 16.8-23.7 25.8-40.5 20.2s-25.8-23.7-20.2-40.5l2.6-7.8c11.4-34.1 36.6-61.9 69.4-76.5l8-3.5c20.8-9.2 43.3-14 66.1-14c44.6 0 84.8 26.8 101.9 67.9L281 232.7l21.4 10.7c15.8 7.9 22.2 27.1 14.3 42.9s-27.1 22.2-42.9 14.3L247 287.3c-10.3-5.2-18.4-13.8-22.8-24.5l-9.6-23-19.3 65.5 49.5 54c5.4 5.9 9.2 13 11.2 20.8l23 92.1c4.3 17.1-6.1 34.5-23.3 38.8s-34.5-6.1-38.8-23.3l-22-88.1-70.7-77.1c-14.8-16.1-20.3-38.6-14.7-59.7l16.9-63.5zM68.7 398l25-62.4c2.1 3 4.5 5.8 7 8.6l40.7 44.4-14.5 36.2c-2.4 6-6 11.5-10.6 16.1L54.6 502.6c-12.5 12.5-32.8 12.5-45.3 0s-12.5-32.8 0-45.3L68.7 398z"/></svg>
        {/if}
      </div>
      <div>
        {#if leg.mode == "BUS"}
          {leg.routeShortName}
        {:else}
          {Math.round(leg.duration/60)}
        {/if}
      </div>
    </div>
  {/each}
</div>