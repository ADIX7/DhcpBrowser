import { injectable } from "tsyringe";

@injectable()
export class BackendApiSettings {
    baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl.endsWith('/') ? baseUrl : baseUrl + '/';
    }
}
