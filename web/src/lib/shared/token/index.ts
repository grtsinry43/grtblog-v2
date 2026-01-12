export function setToken(token: string) {
	localStorage.setItem('authToken', token);
}

export function getToken() {
	return localStorage.getItem('authToken');
}