<script lang="ts">
	import ItineraryLegs from '$lib/ItineraryLegs.svelte';
  import {timestampToString, secondsToDurationString} from '$lib/utils';
	import { fromJSON } from 'postcss';

  export let itinerary;
  export let index;

  let firstTransitLegIndex = itinerary.legs.findIndex(leg => leg.transitLeg);
</script>

<div class="flex flex-col px-4 py-2 w-full rounded-md  border-gray-300 border shadow-sm">
  <div class="flex flex-row text-lg justify-between">
    <span>{timestampToString(itinerary.startTime).time} - {timestampToString(itinerary.endTime).time}</span>
    <span>{secondsToDurationString(itinerary.duration)}</span>
  </div>
  <div class="w-full">
    <ItineraryLegs itinerary={itinerary}/>
  </div>
  {#if firstTransitLegIndex !== -1}
    <span class="text-sm text-gray-500 mt-1 text-left">
      Leaves at <span class="text-black">{timestampToString(itinerary.legs[firstTransitLegIndex].from.departure).time}</span> from <span class="text-black">{itinerary.legs[firstTransitLegIndex].from.name}</span></span>
  {/if}
</div>