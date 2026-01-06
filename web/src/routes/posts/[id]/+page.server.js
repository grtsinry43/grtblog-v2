export const load = async ({ params }) => {
    // MOCK DATA FOR DEMO
    // In a real app, this would fetch from the backend
    const article = {
        title: "Deep Dive into Svelte 5 Runes",
        content: "This is a demo article content to verify the navbar title functionality...",
        date: "2024-03-20"
    };

    return {
        article
    };
};