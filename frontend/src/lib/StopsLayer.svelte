<!-- https://otp.nat.omg.lol/otp/routers/default/vectorTiles/stops/15/16178/10568.pbf -->
<script lang="ts">
  import L from 'leaflet?client';
  import 'leaflet.vectorgrid?client';
  // require('leaflet.vectorgrid');
  import {createEventDispatcher, getContext, setContext, onDestroy, onMount} from 'svelte';
  import {otpBase} from '$lib/utils';

  // export let latLngs;
  // export let options;

  const dispatch = createEventDispatcher();
  
	// let layerPane = pane || getContext('pane');
  let layerGroup = getContext<() => L.LayerGroup>('layerGroup')();
  // let featureGroup = getContext<() => L.FeatureGroup>('featureGroup')();
  let iconSvg: HTMLElement;

  let stops;

  onMount(() => {
    let vectorTileStyling = {
      stops: {
        // icon: L.divIcon({
        //   html: iconSvg.innerHTML,
        //   iconSize: [16, 22],
        //   iconAnchor: [8, 22],
        // })
        icon: L.icon({
          iconUrl: `data:image/svg+xml;base64,${btoa(iconSvg.innerHTML)}`,
          iconSize: [16, 22],
          iconAnchor: [8, 22],
        })
        // icon: new L.Icon.Default(),
      },
    };

    let layerOptions = {
        rendererFactory: L.canvas.tile, // raster, but vector (svg) doen't work properly.
        // attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors, &copy; <a href="https://www.mapbox.com/about/maps/">MapBox</a>',
        vectorTileLayerStyles: vectorTileStyling,
        minNativeZoom: 14,
        maxNativeZoom: 20,
        minZoom: 15
      };

    stops = L.vectorGrid.protobuf("https://otp.nat.omg.lol/otp/routers/default/vectorTiles/stops/{z}/{x}/{y}.pbf", layerOptions)
      .on('click', (e) => dispatch('click', e))
      // .on('mouseover', (e) => dispatch('mouseover', e))
      // .on('mouseout', (e) => dispatch('mouseout', e))
      .addTo(layerGroup);
  });

  setContext('layer', () => stops);

  onDestroy(() => {
    stops.remove();
    stops = undefined;
  });
</script>

<slot />

<div class="hidden">
  <div bind:this={iconSvg}>
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 22"><g clip-path="url(#a)"><circle cx="8" cy="8" r="8" fill="#fff"/><rect width="2" height="8" x="7" y="14" fill="#2B2B2B" rx="1"/><circle cx="8" cy="8" r="7.5" fill="#6366F1"/><path fill="#fff" d="M8 3.5c2.36 0 3.94.62 3.94 1.4v.86c.3 0 .56.25.56.56v1.13c0 .3-.25.56-.56.56v2.81c0 .31-.25.56-.56.56v.57c0 .3-.26.56-.57.56h-.56a.56.56 0 0 1-.56-.56v-.56H6.3v.56c0 .3-.25.56-.56.56h-.56a.56.56 0 0 1-.57-.56v-.56a.56.56 0 0 1-.56-.57V8a.56.56 0 0 1-.56-.56V6.3c0-.3.25-.56.56-.56v-.83C4.06 4.1 5.64 3.5 8 3.5ZM5.19 6.31V8c0 .31.25.56.56.56h1.97V5.75H5.75a.56.56 0 0 0-.56.56Zm3.1 2.25h1.96c.31 0 .56-.25.56-.56V6.31a.56.56 0 0 0-.56-.56H8.28v2.81Zm-2.82 1.97a.56.56 0 1 0 0-1.12.56.56 0 0 0 0 1.12Zm5.06 0a.56.56 0 1 0 0-1.12.56.56 0 0 0 0 1.12ZM9.7 4.91a.28.28 0 0 0-.28-.29H6.59a.28.28 0 0 0-.28.29c0 .15.13.28.28.28h2.82c.15 0 .28-.13.28-.28Z"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h16v22H0z"/></clipPath></defs></svg>
  </div>
</div>