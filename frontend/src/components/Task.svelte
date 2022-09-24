<script>
  import { afterUpdate, createEventDispatcher } from "svelte";

  export let task;

  let editing = false;

  function toggle(e) {
    fetch(`http://localhost:9000/tasks/${task.id}`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ done: e.currentTarget.checked }),
    });
  }

  let input;

  afterUpdate(() => {
    if (editing) input.focus();
  });

  function update() {
    fetch(`http://localhost:9000/tasks/${task.id}`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name: input.value }),
    }).then(() => {
      editing = false;
    });
  }

  const dispatch = createEventDispatcher();

  function handleDelete() {
    fetch(`http://localhost:9000/tasks/${task.id}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    }).then(() => {
      dispatch("delete", {
        id: task.id,
      });
    });
  }
</script>

<div
  class="p-2.5 rounded-lg flex justify-between items-center hover:bg-gray-100"
>
  <div class="flex items-center space-x-4">
    <input
      type="checkbox"
      bind:checked={task.done}
      on:change={toggle}
      value={task.done ? "true" : "false"}
      class="rounded-sm text-blue-800 focus:ring-blue-400"
    />

    {#if editing}
      <form class="relative" on:submit|preventDefault={update}>
        <input
          type="text"
          bind:value={task.name}
          bind:this={input}
          class="h-8 border-gray-300 hover:border-gray-500 outline:none focus:ring-2 focus:ring-blue-300 focus:border-blue-500 rounded-lg"
        />

        <button class="absolute top-1/2 -translate-y-1/2 right-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="w-5 h-5"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M4.5 12.75l6 6 9-13.5"
            />
          </svg>
        </button>
      </form>
    {:else}
      <p class="pointer-events-none">{task.name}</p>
    {/if}
  </div>

  <div class="flex items-center space-x-2">
    <button
      class="text-gray-700 hover:text-blue-800"
      on:click={() => (editing = true)}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-5 h-5"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L6.832 19.82a4.5 4.5 0 01-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 011.13-1.897L16.863 4.487zm0 0L19.5 7.125"
        />
      </svg>

      <span class="sr-only">Edit Task</span>
    </button>

    <button class="text-red-700 hover:text-red-800" on:click={handleDelete}>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-5 h-5"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"
        />
      </svg>

      <span class="sr-only">Delete Task</span>
    </button>
  </div>
</div>
