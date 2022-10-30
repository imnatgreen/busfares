<script lang="ts">
  import { onMount } from 'svelte';
  
  import L from 'leaflet?client';
  import 'leaflet/dist/leaflet.css';

  import { decode } from '@googlemaps/polyline-codec';

  let map;
  let mapContainer;

  let tripPlanPlaceholder = '';
  
  onMount(() => {
    const initialView = new L.LatLng(53.71,-2.24);
    map = L.map(mapContainer, {preferCanvas: false}).setView(initialView, 10);
    let layer = L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}' + (L.Browser.retina ? '@2x.png' : '.png'), {
      attribution:'&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>, &copy; <a href="https://carto.com/attributions">CARTO</a>',
      subdomains: 'abcd',
      maxZoom: 20,
      minZoom: 0
    });
    layer.addTo(map);

    return map.remove;
  });


  // generate a path-dependent if we have VITE_API_URL defined (dev mode) or nor
  const apiUrl = (path: string) => `${import.meta.env.VITE_API_URL || ''}${path}`;
  const otpBase = 'https://otp.nat.omg.lol/otp';
  let from = '53.70346,-2.24653';
  let to = '53.78967,-2.24336';

  
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

  const getPlan = async () => {
    tripPlanPlaceholder = 'getting trip plan...';
    let now = new Date();
    const res = await fetch(otpBase+'/routers/default/plan?'+ new URLSearchParams({
      fromPlace: from,
      toPlace: to,
      date: now.toISOString(),
      mode: 'TRANSIT,WALK',
      maxWalkDistance: '2500',
      arriveBy: 'false',
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
    legs.map((leg, i) => {
      if (leg.legGeometry.points) {
        let mode = leg.mode;
        let lineStyle = {
          color: mode=='WALK' ? 'grey' : '#00839E',
          dashArray: mode=='WALK' ? '.1 11' : null,
          weight: 5,
        };
        let haloLineStyle = {
          color: '#ffffff',
          dashArray: mode=='WALK' ? '.1 11' : null,
          weight: 6,
        };

        let haloLine: L.Polyline = new L.Polyline(decode(leg.legGeometry.points), haloLineStyle).addTo(map);
        let line: L.Polyline = new L.Polyline(decode(leg.legGeometry.points), lineStyle).addTo(map);
        if (i == 0) {
          map.fitBounds(line.getBounds());
        }
      }
    });
  }
</script>

<p><i>plan a journey...</i></p>
<p>from: <input bind:value={from} label="from"/></p>
<p>to: <input bind:value={to} label="to"/></p>
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
<div id="map" bind:this={mapContainer} style="height:600px;width:100%;"/>