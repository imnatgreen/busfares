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