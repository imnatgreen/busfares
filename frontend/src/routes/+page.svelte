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
  let now = new Date();
  let date = now.toISOString().split('T')[0];
  let time = now.getHours().toString().padStart(2, '0') + ':' + now.getMinutes();
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
  let currentItinerary = 0;

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

  $: if(tripPlan) {
    drawLine(tripPlan.plan.itineraries[currentItinerary].legs);
  }

  const drawLine = (legs) => {
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
<svelte:window on:resize={resizeMap} />

<p><i>plan a journey...</i></p>
<p>from: <input bind:value={from} label="from"/></p>
<p>to: <input bind:value={to} label="to"/></p>
<p>date: <input type="date" bind:value={date} label="date"/> time: <input type="time" bind:value={time} label="time"/></p>
<p>arrive by? <input type="checkbox" bind:checked={arriveBy} label="arrive by"/></p>
<button on:click={getPlan}>find journey</button>
<hr>
{#if tripPlan}
  {#each tripPlan.plan.itineraries as itinerary, i}
    <label><input type=radio bind:group={currentItinerary} name="currentItinerary" value={i}>itinerary {i+1}</label>
    {#each itinerary.legs as leg}
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
  <p>{tripPlanPlaceholder}</p>
{/if}
<div style="height:600px;width:100%;">
  <Leaflet bind:map view={initialView} zoom={10} on:click={mapClick}>
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
  </Leaflet>
</div>

<div style="display: none;"> 
  <div bind:this={mapClickPopupContent}>
    {#if mapClickPopupEvent}
    <p>{mapClickPopupEvent ? latLngString(mapClickPopupEvent.latlng) : ""}</p>
    <button on:click={() => setLocation('from', mapClickPopupEvent.latlng)}>from here</button>
    <button on:click={() => setLocation('to', mapClickPopupEvent.latlng)}>to here</button>
    {/if}
  </div>
</div>