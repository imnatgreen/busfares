// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
// and what to do when importing types
declare namespace App {
	// interface Locals {}
	// interface PageData {}
	// interface Error {}
	// interface Platform {}
}

interface ImportMetaEnv {
	VITE_API_URL: string;
}

declare module 'leaflet?client' {
	import all from 'leaflet';
	export = all;
}

declare module 'leaflet.vectorgrid?client' {
	import all from 'leaflet.vectorgrid';
	export = all;
}

declare module 'ol-mapbox-style?client' {
	import all from 'ol-mapbox-style';
	export = all;
}