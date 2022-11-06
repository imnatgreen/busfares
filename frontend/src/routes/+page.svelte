<script lang="ts">
  import { onMount } from 'svelte';
  
  import L from 'leaflet?client';
  import 'leaflet/dist/leaflet.css';

  import { decode } from '@googlemaps/polyline-codec';
	import Leaflet from '$lib/Leaflet.svelte';
	import Polyline from '$lib/Polyline.svelte';
	import Popup from '$lib/Popup.svelte';

  let map;
  // let mapContainer;

  let tripPlanPlaceholder = '';
  
  const initialView = [53.71,-2.24];
  function resizeMap() {
	  if(map) { map.invalidateSize(); }
  }
	
	function resetMapView() {
		map.setView(initialView, 10);
	}

  // generate a path-dependent if we have VITE_API_URL defined (dev mode) or nor
  const apiUrl = (path: string) => `${import.meta.env.VITE_API_URL || ''}${path}`;
  const otpBase = 'https://otp.nat.omg.lol/otp';
  let from = '53.70346,-2.24653';
  let to = '53.78967,-2.24336';
  let now = new Date()
  let date = now.toISOString().split('T')[0];
  let time = now.getHours() + ':' + now.getMinutes();
  let arriveBy = false;
  
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
  $: mapLines = [];

  const getPlan = async () => {
    tripPlan = undefined;
    tripPlanPlaceholder = 'getting trip plan...';
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
    let json = await res.json();
    json = await getFares(json);
    tripPlan = json;
    testDrawLine(tripPlan);
    
  }

  const testDrawLine = (tripPlan) => {
    let legs = tripPlan.plan.itineraries[0].legs;
    mapLines = [];
    legs.map((leg, i) => {
      if (leg.legGeometry.points) {
        let line = {
          id: i,
          latLngs: decode(leg.legGeometry.points),
          mode: leg.mode,
        }
        mapLines = [...mapLines, line];
      }
    });
  }
</script>
<svelte:window on:resize={resizeMap} />

{@debug tripPlan}

<p><i>plan a journey...</i></p>
<p>from: <input bind:value={from} label="from"/></p>
<p>to: <input bind:value={to} label="to"/></p>
<p>date: <input type="date" bind:value={date} label="date"/> time: <input type="time" bind:value={time} label="time"/></p>
<p>arrive by? <input type="checkbox" bind:checked={arriveBy} label="arrive by"/></p>
<button on:click={getPlan}>find journey</button>
<hr>
{#if tripPlan}
  {#each tripPlan.plan.itineraries[0].legs as leg}
    <p>{new Date(leg.startTime).toLocaleString()}: {leg.mode} from {leg.from.name} to {leg.to.name}</p>
    {#if leg.fares}
      <select>
        {#each leg.fares as fare}
          <option value={fare.salesOfferPackage.id}>{fare.preassignedFareProduct.name}: {fare.amount.currency} {fare.amount.number}</option>
        {/each}
      </select>
    {/if}
  {/each}
{:else}
  <p>{tripPlanPlaceholder}</p>
{/if}
<div style="height:600px;width:100%;">
  <Leaflet bind:map view={initialView} zoom={10}>
    {#each mapLines as line}
      <Polyline latLngs={line.latLngs} options={{
        color: '#ffffff',
        dashArray: line.mode=='WALK' ? '.1 11' : null,
        weight: 6,
      }}/>
      <Polyline latLngs={line.latLngs} options={{
        color: line.mode=='WALK' ? 'grey' : '#00839e',
        dashArray: line.mode=='WALK' ? '.1 11' : null,
        weight: 5,
      }}>
        <Popup>{line.mode}</Popup>
      </Polyline>
    {/each}
    <Popup>test</Popup>
  </Leaflet>
</div>