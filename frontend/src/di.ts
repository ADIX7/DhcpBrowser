import { Lifecycle, container } from "tsyringe";
import { BackendApiSettings } from "./configuration/backend-api-settings";
import { BackendApi } from "./services/backend-api";

console.log(import.meta.env);

container.register(BackendApiSettings, { useValue: new BackendApiSettings((window as any).BACKEND_API_URL ?? import.meta.env.VITE_BACKEND_API_URL ?? window.location.origin) });
container.register(BackendApi, { useClass: BackendApi }, { lifecycle: Lifecycle.Singleton });
