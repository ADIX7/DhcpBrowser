import { createSignal, type Component, onMount, For } from 'solid-js';

import styles from './App.module.scss';
import { container } from 'tsyringe';
import { BackendApi } from './services/backend-api';

const App: Component = () => {

    const backendApi = container.resolve(BackendApi);
    const [ipv4Leases, setIpv4Leases] = createSignal<Ipv4Lease[]>([]);
    const [loading, setLoading] = createSignal(true);

    onMount(async () => {
        const leases = await backendApi.getIpv4Leases();
        setIpv4Leases(leases);
        setLoading(false);
    });

    return (
        <div class={styles.lease_container}>
            <For each={ipv4Leases()}>
                {lease => (
                    <div>
                        <span>{lease.ipAddress}</span>
                        <span>{lease.hwAddress}</span>
                        <span>{new Date(lease.expiresAt * 1000).toString()}</span>
                    </div>
                )}
            </For>
        </div>
    );
};

export default App;
