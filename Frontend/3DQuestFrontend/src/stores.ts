import { writable, type Writable } from "svelte/store";

export const isLogged : Writable<boolean> = writable<boolean>(false);