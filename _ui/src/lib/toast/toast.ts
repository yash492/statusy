import toast from 'svelte-french-toast';

export class Toast {
	private style = 'border-radius: 50px; background: #333; color: #fff;';

	success(msg: string) {
		return toast.success(msg, { style: this.style });
	}

	error(msg: string) {
		return toast.error(msg, { style: this.style });
	}

	custom(icon: string, msg: string) {
		return toast(msg, {
			icon: icon,
			style: this.style
		});
	}
}
