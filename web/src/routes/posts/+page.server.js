import { fail } from '@sveltejs/kit';

export const load = async ({ fetch }) => {
	const response = await fetch('http://localhost:8080/article/all?page=1&pageSize=10');
	const result = await response.json();	const articles = result.data;
	return {
		articles,
	};
};