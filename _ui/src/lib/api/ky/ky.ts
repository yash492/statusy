import { PUBLIC_API_SERVER_ROUTE } from "$env/static/public";
import ky from "ky";

const KyClient = ky.create({
    prefixUrl: PUBLIC_API_SERVER_ROUTE,
    headers: {
        accept: "application/json",
    },
});

export default KyClient;