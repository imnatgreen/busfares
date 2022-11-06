<script lang="ts">
  import L from 'leaflet?client';
  import {createEventDispatcher, getContext, setContext, onDestroy} from 'svelte';

  export let latLngs;
  export let options;

  const dispatch = createEventDispatcher();
  
	// let layerPane = pane || getContext('pane');
  let layerGroup = getContext<() => L.LayerGroup>('layerGroup')();

  export let line: L.Polyline = new L.Polyline(latLngs, options)
    .on('click', (e) => dispatch('click', e))
    .on('mouseover', (e) => dispatch('mouseover', e))
    .on('mouseout', (e) => dispatch('mouseout', e))
    .addTo(layerGroup);

  setContext('layer', () => line);

  onDestroy(() => {
    line.remove();
    line = undefined;
  });

  $: line.setStyle(options);

  $: {
    line.setLatLngs(latLngs);
    line.redraw();
  }
</script>

<slot />