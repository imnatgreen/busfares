<script lang="ts">
	import ItineraryLegs from '$lib/ItineraryLegs.svelte';
  import {timestampToString, secondsToDurationString} from '$lib/utils';

  import { fly } from 'svelte/transition'

  export let itinerary;
  export let totalFare;
  // export let index;

  let firstTransitLegIndex = itinerary.legs.findIndex(leg => leg.transitLeg);
</script>

<div class="flex flex-col px-4 py-2 w-full rounded-md  border-gray-300 hover:bg-indigo-50 hover:bg-opacity-30 border shadow-sm hover:shadow-md transition-all">
  <div class="flex flex-row text-lg justify-between">
    <span>{timestampToString(itinerary.startTime).time} - {timestampToString(itinerary.endTime).time}</span>
    <span>{secondsToDurationString(itinerary.duration)}</span>
  </div>
  <div class="w-full">
    <ItineraryLegs itinerary={itinerary}/>
  </div>
  {#if firstTransitLegIndex !== -1}
  <div class="flex flex-row text-lg justify-between align-middle">
    <span class="text-sm text-gray-500 mt-1 text-left">
      Leaves at <span class="text-black">{timestampToString(itinerary.legs[firstTransitLegIndex].from.departure).time}</span> from <span class="text-black">{itinerary.legs[firstTransitLegIndex].from.name}</span>
    </span>
    <div class="text-sm text-gray-500 mt-1 ml-2 text-left">
      {#if totalFare}
        <span transition:fly="{{x:25}}">{new Intl.NumberFormat('en-GB',{ style: 'currency', currency: 'GBP' }).format(totalFare)}</span>
      {/if}
      </div>
  </div>
  {/if}
</div>