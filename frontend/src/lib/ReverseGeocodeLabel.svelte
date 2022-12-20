<script lang="ts">
  import {reverseGeocode, latLngString, placeToName} from '$lib/utils';

  export let latlng;
  export let place = null;
  let label = latLngString(latlng);

  $: (async () => {
    place = await reverseGeocode(latlng.lat, latlng.lng);
    label = place ? placeToName(place) : latLngString(latlng);
  })();
</script>

<span>{label}</span>