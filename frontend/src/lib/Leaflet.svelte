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
	let mapFeatureGroup;
  $: mapProp = map;
	
  export const getMap = () => map;
  export const getMapFeatureGroup = () => mapFeatureGroup;
  setContext('layerGroup', getMap);
  setContext('featureGroup', getMapFeatureGroup);
  setContext('layer', getMap);
  setContext('map', getMap);
	
  function createLeaflet(node) {
    map = L.map(node, {zoomControl:false})
      .addControl(L.control.scale({metric: true, imperial: false, position:'bottomright'}))
      .addControl(L.control.zoom({position:'bottomright'}))
      .on('zoom', (e) => dispatch('zoom', e))
      .on('click', (e) => dispatch('click', e));
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
		
		mapFeatureGroup = L.featureGroup().addTo(map);
		
    return {
      destroy() {
        map.remove();
        map = undefined;
      },
    };
  }
	
	export function fitMapToLines() {
		if(mapFeatureGroup) {
			if(mapFeatureGroup.getLayers().length > 0) {
				map.flyToBounds(mapFeatureGroup.getBounds().pad(0.15), {duration: 0.2});
			}
		}
	}
		
	$: if(map) {
		if(bounds) {
      map.fitBounds(bounds)
		} else {
			map.setView(view, zoom);
		}
	}
</script> 

<div class="{classes}" use:createLeaflet>
  {#if map}
    <slot {map} />
  {/if}
</div>

<style>
  :global(.leaflet-control-container) {
    @apply static;
  }
  :global(.leaflet-container .leaflet-control-attribution) {
	height: 22px;
	opacity: 0.7;
	border-radius: 11px;
	background-color: #fff;
	margin: 10px;
	position: relative;
	top: -5px;
}
:global(.leaflet-container .leaflet-popup-close-button) {
	@apply text-indigo-500;
	font-size: 25px;
	font-weight: 100;
	width: 32px;
	height: 32px;
	top: 12px;
	right: 4px;
	padding: 0;
}
:global(.leaflet-container .leaflet-control-attribution a) {
	height: 11px;
	@apply font-sans;
	font-size: 11px;
	font-weight: normal;
	font-stretch: normal;
	font-style: normal;
	line-height: normal;
	letter-spacing: -0.7px;
	text-align: right;
	color: #666;
	vertical-align: -3px;
}
:global(div.leaflet-bottom.leaflet-right div.leaflet-control-zoom) {
	width: 36px;
	height: 73px;
	border-radius: 5px;
	background-color: #fff;
	margin-right: 30px;
	margin-bottom: 24px;
	box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}
:global(div.leaflet-bottom.leaflet-right div.leaflet-control-zoom.leaflet-bar) {
	border: none;
}
:global(div.leaflet-bottom.leaflet-right div.leaflet-control-zoom .icon) {
	font-size: 18px;
}
:global(div.leaflet-bottom.leaflet-right div.leaflet-control-zoom a) {
	display: flex;
	align-items: center;
	justify-content: space-around;
	width: 36px;
	height: 36px;
	line-height: 18px;
	font-size: 18px;
	@apply text-indigo-400;
}
/* Fix to default leaflet behavior */
:global(.leaflet-map-pane svg) {
	position: relative;
}
:global(.leaflet-map-pane svg.icon-badge) {
	transform: translate(-0.5em, -2.5em);
	border-radius: 50%;
}
:global(.leaflet-map-pane svg.icon-badge > .badge-circle) {
	stroke-width: 14%;
}
:global(div.leaflet-control-scale.leaflet-control) {
	margin-right: 30px;
	margin-bottom: 20px;
	cursor: grab;
}
:global(div.leaflet-control-scale-line) {
	text-align: right;
	margin-bottom: -4px;
	cursor: grab;
	background: rgba(255, 255, 255, 0);
	border: none;
}
:global(.leaflet-control-scale::after) {
	content: '';
	display: block;
	border-bottom: 1px solid #888;
	border-left: 1px solid #888;
	border-right: 1px solid #888;
	height: 4px;
	background: none;
	cursor: grab;
}

</style>