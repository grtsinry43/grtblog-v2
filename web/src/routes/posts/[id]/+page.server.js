export const load = async ({ params, fetch }) => {
    const { id } = params;
    const shortUrl = id; // Assuming `id` in the route actually corresponds to `shortUrl`
    const response = await fetch(`http://localhost:8080/article/${shortUrl}`);
    const result = await response.json();
    const article = result.data;
    return {
        article
    };
};