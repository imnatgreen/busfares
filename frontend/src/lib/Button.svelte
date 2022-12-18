<script lang="ts">
	import { createEventDispatcher } from "svelte";
  
  export let classes = '';
  export let wrapClasses = '';
  
  const dispatch = createEventDispatcher();

  function handleClick() {
    dispatch('click');
  }
</script>

<div class="relative {wrapClasses}">
  {#if $$slots.icon}
    <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-2 text-gray-500">
      <slot name="icon"></slot>
    </div>
  {/if}
  <button on:click={handleClick} class="form-input px-2 py-1 {$$slots.icon && !classes.includes('pl-') ? 'pl-8' : ''} block rounded-md border-gray-300 shadow-sm
  focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 hover:bg-indigo-50 transition-all {classes}">
    <slot></slot>
  </button>
</div>