<script lang="ts">
  import { onMount } from 'svelte';
  import { fade, fly } from 'svelte/transition'
  import { otpBase, apiUrl, getLineColour } from '$lib/utils';
  
  import L from 'leaflet?client';
  import 'leaflet/dist/leaflet.css';

  import { decode } from '@googlemaps/polyline-codec';
	import Leaflet from '$lib/Leaflet.svelte';
	import Polyline from '$lib/Polyline.svelte';
	import Popup from '$lib/Popup.svelte';
	import ItinerarySummary from '$lib/ItinerarySummary.svelte';

  let map;

  let tripPlanPlaceholder = '';
  let gettingPlan = false;
  
  const initialView = [53.71,-2.24];
  function resizeMap() {
	  if(map) { map.invalidateSize(); }
  }
	
	function resetMapView() {
		map.setView(initialView, 10);
	}

  // generate a path-dependent if we have VITE_API_URL defined (dev mode) or nor

  let from = '53.70346,-2.24653';
  let to = '53.78967,-2.24336';
  let now = new Date();
  let datePicker;
  let timePicker;
  let date = now.toISOString().split('T')[0];
  let time = now.getHours().toString().padStart(2, '0') + ':' + now.getMinutes().toString().padStart(2, '0');
  let arriveBy = false;
  let editSearch = true;
  
  const openEditSearch = () => {
    editSearch = true;
  }

  const getFares = async (tripPlan: object) => {
    tripPlanPlaceholder = 'getting fares...'
    const res = await fetch(apiUrl('/api/getfares'), {
      method: 'POST',
      // headers: {
      //   'Content-Type': 'application/json',
      // },
      body: JSON.stringify(tripPlan),
    });
    return await res.json();
  };

  let tripPlan;
  let tripPlanFares;
  $: mapLines = [];
  let currentItinerary = 0;

  const getPlan = async () => {
    tripPlan = undefined;
    tripPlanFares = undefined;
    tripPlanPlaceholder = 'getting trip plan...';
    gettingPlan = true;
    editSearch = false;
    const res = await fetch(otpBase+'/routers/default/plan?'+ new URLSearchParams({
      fromPlace: from,
      toPlace: to,
      date: date,
      time: time,
      mode: 'TRANSIT,WALK',
      maxWalkDistance: '2500',
      arriveBy: arriveBy ? 'true' : 'false',
      wheelchair: 'false',
      showIntermediateStops: 'true',
      locale: 'en'
    }));
    tripPlan = await res.json();
    gettingPlan = false;
    tripPlanFares = await getFares(tripPlan);
    tripPlanFares.forEach((legs, i) => {
      legs.forEach((fares, l) => {
        tripPlan.plan.itineraries[i].legs[l].fares = fares;
      })
    });
    for (const [i, itinerary] of tripPlan.plan.itineraries.entries()) {
      for (const [l, leg] of itinerary.legs.entries()) {
        if (leg.mode === 'BUS') {
          let col = await getLineColour(leg.agencyId.substring(2), leg.routeShortName);
          if (col != '') {
            tripPlan.plan.itineraries[i].legs[l].colour = col;
          }
        }
      }
    }
  }

  $: if(tripPlan) {
    if (tripPlan.plan.itineraries.length > 0) {
      drawLine(tripPlan.plan.itineraries[currentItinerary].legs);
    }
  }

  const drawLine = (legs) => {
    mapLines = [];
    legs.map((leg, i) => {
      if (leg.legGeometry.points) {
        let line = {
          id: i,
          latLngs: decode(leg.legGeometry.points),
          mode: leg.mode,
          colour: leg.colour ? leg.colour : '#6366F1'
        }
        mapLines = [...mapLines, line];
      }
    });
  }

  let mapClickPopupContent;
  let mapClickPopupEvent;
  let mapClickPopup;
  onMount(async () => {
    mapClickPopup = L.popup();
  })

  const mapClick = (e) => {
    mapClickPopupEvent = e.detail;
    mapClickPopup
      .setLatLng(mapClickPopupEvent.latlng)
      .setContent(mapClickPopupContent)
      .openOn(map);
    console.log(mapClickPopupEvent.latlng);      
  }

  const latLngString = (latLng) => {
    return latLng.lat.toFixed(5) + ',' + latLng.lng.toFixed(5);
  }

  const setLocation = (location, latlng) => {
    editSearch = true;
    if (location == 'from') {
      from = latLngString(latlng);
    } else if (location == 'to') {
      to = latLngString(latlng);
    }
    if (mapClickPopup) {
      mapClickPopup.remove();
    }
  }
</script>
<svelte:head>
  <title>bus fares.</title>
</svelte:head>
<svelte:window on:resize={resizeMap} />

<div class="flex flex-col md:flex-row w-full h-[calc(100vh-env(safe-area-inset-bottom))] p-4 gap-4 bg-gradient-to-br from-pink-100 via-purple-100 to-indigo-200">
  <div class="h-1/2 md:h-full md:w-1/2 lg:w-2/5 overflow-hidden flex-none rounded-lg shadow-lg bg-white">
    <div class="w-full h-full p-4 overflow-auto">
      <div class="grid grid-rows-[1fr] grid-cols-[1fr]">
        {#if editSearch}
          <div transition:fly="{{y:-100, duration:500}}" class="row-[1] col-[1]">
            <h1 class="text-2xl">bus fares.</h1>
            <div class="text-gray-500 text-sm">
              <p>get bus, but cheap.</p>
              <div class="flex flex-row space-x-px items-center">
                <span>made with</span>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-5 w-5 pb-px hover:text-red-400 transition"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fill-rule="evenodd"
                    d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z"
                    clip-rule="evenodd"
                  />
                </svg>
                <span>by nathan :)</span>
              </div>
            </div>
            <p class="mb-1 mt-4 text-gray-500 text-sm">okay, first, where to?</p>
            <div class="relative">
              <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                <span class="text-gray-500">from</span>
              </div>
              <input
                aria-label="Journey to location"
                type="text"
                class="
                  mt-1 px-2 pl-14 py-1 block w-full rounded-md  border-gray-300 shadow-sm
                  focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50
                "
                bind:value={from} name="from" label="from"
              />
            </div>
            <div class="relative">
              <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                <span class="text-gray-500">to</span>
              </div>
              <input
                aria-label="Journey to location"
                type="text"
                class="
                  mt-1 px-2 pl-14 py-1 block w-full rounded-md  border-gray-300 shadow-sm
                  focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50
                "
                bind:value={to} name="to" label="to"
              />
            </div>

            <p class="mb-1 mt-4 text-gray-500 text-sm">great! now, when?</p>
            <div class="mt-1 flex flex-row gap-1 flex-wrap">
              <div class="relative">
                <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-2 text-gray-500">
                  {#if arriveBy}
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                      <path fill-rule="evenodd" d="M3 4.25A2.25 2.25 0 015.25 2h5.5A2.25 2.25 0 0113 4.25v2a.75.75 0 01-1.5 0v-2a.75.75 0 00-.75-.75h-5.5a.75.75 0 00-.75.75v11.5c0 .414.336.75.75.75h5.5a.75.75 0 00.75-.75v-2a.75.75 0 011.5 0v2A2.25 2.25 0 0110.75 18h-5.5A2.25 2.25 0 013 15.75V4.25z" clip-rule="evenodd" />
                      <path fill-rule="evenodd" d="M19 10a.75.75 0 00-.75-.75H8.704l1.048-.943a.75.75 0 10-1.004-1.114l-2.5 2.25a.75.75 0 000 1.114l2.5 2.25a.75.75 0 101.004-1.114l-1.048-.943h9.546A.75.75 0 0019 10z" clip-rule="evenodd" />
                    </svg>
                  {:else}
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                      <path fill-rule="evenodd" d="M3 4.25A2.25 2.25 0 015.25 2h5.5A2.25 2.25 0 0113 4.25v2a.75.75 0 01-1.5 0v-2a.75.75 0 00-.75-.75h-5.5a.75.75 0 00-.75.75v11.5c0 .414.336.75.75.75h5.5a.75.75 0 00.75-.75v-2a.75.75 0 011.5 0v2A2.25 2.25 0 0110.75 18h-5.5A2.25 2.25 0 013 15.75V4.25z" clip-rule="evenodd" />
                      <path fill-rule="evenodd" d="M6 10a.75.75 0 01.75-.75h9.546l-1.048-.943a.75.75 0 111.004-1.114l2.5 2.25a.75.75 0 010 1.114l-2.5 2.25a.75.75 0 11-1.004-1.114l1.048-.943H6.75A.75.75 0 016 10z" clip-rule="evenodd" />
                    </svg>
                  {/if}
                </div>
                <select
                  aria-label="Leave at time or arrive by time?"
                  bind:value={arriveBy} label="arrive by"
                  class="
                    px-2 pl-7 {arriveBy ? 'pr-7' : 'pr-6'} py-1 block rounded-md  border-gray-300 shadow-sm
                    focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50
                  "
                >
                  <option value={false}>leave at</option>
                  <option value={true}>arrive by</option>
                </select>
              </div>
              <div class="relative">
                <button aria-label="Open date picker" on:click={() => datePicker.showPicker()} class="absolute inset-y-0 left-0 flex items-center pl-2 text-gray-500">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                    <path d="M5.25 12a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H6a.75.75 0 01-.75-.75V12zM6 13.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V14a.75.75 0 00-.75-.75H6zM7.25 12a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H8a.75.75 0 01-.75-.75V12zM8 13.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V14a.75.75 0 00-.75-.75H8zM9.25 10a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H10a.75.75 0 01-.75-.75V10zM10 11.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V12a.75.75 0 00-.75-.75H10zM9.25 14a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H10a.75.75 0 01-.75-.75V14zM12 9.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V10a.75.75 0 00-.75-.75H12zM11.25 12a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H12a.75.75 0 01-.75-.75V12zM12 13.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V14a.75.75 0 00-.75-.75H12zM13.25 10a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H14a.75.75 0 01-.75-.75V10zM14 11.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V12a.75.75 0 00-.75-.75H14z" />
                    <path fill-rule="evenodd" d="M5.75 2a.75.75 0 01.75.75V4h7V2.75a.75.75 0 011.5 0V4h.25A2.75 2.75 0 0118 6.75v8.5A2.75 2.75 0 0115.25 18H4.75A2.75 2.75 0 012 15.25v-8.5A2.75 2.75 0 014.75 4H5V2.75A.75.75 0 015.75 2zm-1 5.5c-.69 0-1.25.56-1.25 1.25v6.5c0 .69.56 1.25 1.25 1.25h10.5c.69 0 1.25-.56 1.25-1.25v-6.5c0-.69-.56-1.25-1.25-1.25H4.75z" clip-rule="evenodd" />
                  </svg>
                </button>
                <input
                  aria-label="Date picker"
                  bind:this={datePicker}
                  type="date" bind:value={date} label="date"
                  class="
                    px-2 pl-8 py-1 block rounded-md  border-gray-300 shadow-sm
                    focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50
                  "
                />
              </div>
              <div class="relative">
                <button aria-label="Open time picker" on:click={() => timePicker.showPicker()} class="absolute inset-y-0 left-0 flex items-center pl-2 text-gray-500">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z" clip-rule="evenodd" />
                  </svg>          
                </button>
                <input
                  aria-label="Time picker"
                  bind:this={timePicker}
                  type="time" bind:value={time} label="time"
                  class="
                    px-2 pl-8 py-1 block rounded-md  border-gray-300 shadow-sm
                    focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50
                  "
                />
              </div>
            </div>

            <p class="mb-1 mt-4 text-gray-500 text-sm">and, finally...</p>
            <div class="relative mt-1">
              <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-2 text-gray-500">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                  <path fill-rule="evenodd" d="M8.157 2.175a1.5 1.5 0 00-1.147 0l-4.084 1.69A1.5 1.5 0 002 5.251v10.877a1.5 1.5 0 002.074 1.386l3.51-1.453 4.26 1.763a1.5 1.5 0 001.146 0l4.083-1.69A1.5 1.5 0 0018 14.748V3.873a1.5 1.5 0 00-2.073-1.386l-3.51 1.452-4.26-1.763zM7.58 5a.75.75 0 01.75.75v6.5a.75.75 0 01-1.5 0v-6.5A.75.75 0 017.58 5zm5.59 2.75a.75.75 0 00-1.5 0v6.5a.75.75 0 001.5 0v-6.5z" clip-rule="evenodd" />
                </svg>
              </div>
              <button on:click={getPlan} class="form-input pl-8 px-2 py-1 block rounded-md border-gray-300 shadow-sm
              focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 hover:bg-indigo-50">find journey</button>
            </div>
            <hr class="my-2">
          </div>
        {:else}
          <div transition:fly="{{y:100, duration:500}}" class="row-[1] col-[1] text-gray-500 text-sm">
            <div>journey from <span class="text-black">{from}</span> to <span class="text-black">{to}</span>, {arriveBy ? 'arriving by' : 'leaving at'} {date} {time}</div>
            <div class="relative mt-1">
              <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-2">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4">
                  <path d="M2.695 14.763l-1.262 3.154a.5.5 0 00.65.65l3.155-1.262a4 4 0 001.343-.885L17.5 5.5a2.121 2.121 0 00-3-3L3.58 13.42a4 4 0 00-.885 1.343z" />
                </svg>
              </div>
              <button on:click={openEditSearch} class="form-input text-sm pl-7 px-2 py-1 block rounded-md border-gray-300 shadow-sm
              focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 hover:bg-indigo-50">edit</button>
            </div>
            <hr class="my-2">
          </div>
        {/if}
      </div>
      {#key editSearch}
        <div in:fly="{
            editSearch?{y:-100,duration:500, delay:0, opacity:1}
            :{y:25,duration:500,delay:500}
          }" class="{editSearch ? 'opacity-100':''} transition-opacity duration-500 -z-10" >
          {#if tripPlan}
            <div in:fly="{
              editSearch?{y:-100,duration:500, delay:0, opacity:1}
              :{y:25,duration:500,delay:500}
            }">
              {#if tripPlan.plan.itineraries.length > 0}
                {#each tripPlan.plan.itineraries as itinerary, i}
                  <ItinerarySummary
                    itinerary={itinerary}
                    index={i}
                    on:select={() => currentItinerary = i}
                  />
                  <label><input type=radio bind:group={currentItinerary} name="currentItinerary" value={i}>itinerary {i+1}</label>
                  {#each itinerary.legs as leg, l}
                    <p>{new Date(leg.startTime).toLocaleString()}: {leg.mode} {leg.mode=='BUS'? '('+leg.routeShortName+' towards '+leg.headsign+')' : ''} from {leg.from.name} to {leg.to.name}</p>
                    {#if leg.fares}
                      <select>
                        {#each leg.fares as fare}
                          <option value={fare.salesOfferPackage.id}>{fare.preassignedFareProduct.name}: {fare.amount.currency} {fare.amount.number}</option>
                        {/each}
                      </select>
                    {/if}
                  {/each}
                {/each}
              {:else}
                <p>no itineraries found</p>
              {/if}
            </div>
          {:else}
            {#if gettingPlan}
              <svg out:fade="{{duration:250}}" class="animate-spin-fast w-8 h-8" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path d="M12,1A11,11,0,1,0,23,12,11,11,0,0,0,12,1Zm0,19a8,8,0,1,1,8-8A8,8,0,0,1,12,20Z" opacity=".25"/>
                <path d="M10.14,1.16a11,11,0,0,0-9,8.92A1.59,1.59,0,0,0,2.46,12,1.52,1.52,0,0,0,4.11,10.7a8,8,0,0,1,6.66-6.61A1.42,1.42,0,0,0,12,2.69h0A1.57,1.57,0,0,0,10.14,1.16Z"/>
              </svg>
            {/if}
            <!-- <p>{tripPlanPlaceholder}</p> -->
          {/if}
        </div>
      {/key}
    </div>
  </div>
  <div class="h-1/2 md:h-full w-full md:w-1/2 lg:w-3/5 md:pr-4 flex-none">
    <Leaflet classes="h-[calc(50vh-2rem)] md:h-[calc(100vh-2rem)] w-full rounded-lg overflow-hidden z-10 shadow-lg" bind:map view={initialView} zoom={10} on:click={mapClick}>
      {#each mapLines as line}
        <Polyline latLngs={line.latLngs} options={{
          color: '#ffffff',
          dashArray: line.mode=='WALK' ? '.1 11' : null,
          weight: 6,
        }}/>
        <Polyline latLngs={line.latLngs} options={{
          color: line.mode=='WALK' ? 'grey' : line.colour,
          dashArray: line.mode=='WALK' ? '.1 11' : null,
          weight: 5,
        }}>
          <Popup>{line.mode}</Popup>
        </Polyline>
      {/each}
    </Leaflet>
  </div>
</div>

<div class="hidden">
  <div bind:this={mapClickPopupContent}>
    {#if mapClickPopupEvent}
      <div class="flex flex-col gap-2 font-sans text-sm">
        <span>{mapClickPopupEvent ? latLngString(mapClickPopupEvent.latlng) : ""}</span>
        <div class="flex gap-2">
          <button on:click={() => setLocation('from', mapClickPopupEvent.latlng)} class="form-input px-2 py-1 block text-sm rounded-md border-gray-300 shadow-sm
            focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 hover:bg-indigo-50">from here</button>
          <button on:click={() => setLocation('to', mapClickPopupEvent.latlng)} class="form-input px-2 py-1 block text-sm rounded-md border-gray-300 shadow-sm
            focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 hover:bg-indigo-50">to here</button>
        </div>
      </div>
    {/if}
  </div>
</div>