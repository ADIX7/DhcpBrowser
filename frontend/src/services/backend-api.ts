import { injectable } from "tsyringe";
import { BackendApiSettings } from "../configuration/backend-api-settings";

@injectable()
export class BackendApi {
    constructor(private readonly settings: BackendApiSettings) { }

    public async getIpv4Leases(): Promise<Ipv4LeasesResult> {
        return fetch(this.settings.baseUrl + 'api/ipv4-leases')
            .then(response => response.json());
    }
}
