// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	interface ViewTransition {
		ready: Promise<void>;
	}

	interface Document {
		startViewTransition?: (callback: () => void) => ViewTransition;
	}

	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
