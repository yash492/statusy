import KyClient from '$lib/ky/ky.js';

export function load({ params }) {
    KyClient.get("/")
}