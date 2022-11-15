<script>
  import { createEventDispatcher, setContext } from 'svelte';
  import L from 'leaflet?client';
  import 'leaflet/dist/leaflet.css';
  export let classes = '';
	
	// Must set either bounds, or view and zoom.
  export let bounds = undefined;
	export let view = undefined;
	export let zoom = undefined;
  let mapProp = undefined;
  export { mapProp as map };

	export const invalidateSize = () => map?.invalidateSize();
  
	const dispatch = createEventDispatcher();

	let map;
  $: mapProp = map;
	
  export const getMap = () => map;
  setContext('layerGroup', getMap);
  setContext('layer', getMap);
  setContext('map', getMap);
	
  function createLeaflet(node) {
    map = L.map(node)
      .on('zoom', (e) => dispatch('zoom', e))
      .on('click', (e) => {dispatch('click', e);
    console.log('clicked'+e.latlng.toString())});
		if(bounds) {
      map.fitBounds(bounds)
		} else {
			map.setView(view, zoom);
		}

      
    L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}' + (L.Browser.retina ? '@2x.png' : '.png'), {
      attribution:'&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>, &copy; <a href="https://carto.com/attributions">CARTO</a>',
      subdomains: 'abcd',
      maxZoom: 20,
      minZoom: 0
    }).addTo(map);
		
    return {
      destroy() {
        map.remove();
        map = undefined;
      },
    };
  }
	
	$: if(map) {
		if(bounds) {
      map.fitBounds(bounds)
		} else {
			map.setView(view, zoom);
		}
	}
</script> 

<style>
  :global(.leaflet-control-container) {
    position: static;
  }
</style>

<div class="{classes}" use:createLeaflet>
  {#if map}
    <slot {map} />
  {/if}
</div>