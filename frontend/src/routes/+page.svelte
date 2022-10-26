<script lang="ts">
  // generate a path-dependent if we have VITE_API_URL defined (dev mode) or nor
  const apiUrl = (path: string) => `${import.meta.env.VITE_API_URL || ''}${path}`;

  const getTestData = async () => {
    const res = await fetch(apiUrl('/api/'));
    return await res.text();
  };
</script>
<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>
{#await getTestData()}
  <p>loading...</p>
{:then data}
  <p>{data}</p>
{/await}