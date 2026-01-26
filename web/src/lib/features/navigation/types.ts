export interface NavMenuItem {
	id: number;
	name: string;
	url: string;
	icon?: string | null;
	children?: NavMenuItem[];
}
