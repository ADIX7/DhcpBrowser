import { createSignal, type Component, onMount, For, onCleanup } from 'solid-js';

import styles from './App.module.scss';
import { container } from 'tsyringe';
import { BackendApi } from './services/backend-api';
import { IpLeaseBlock } from './IpLeaseBlock';

const App: Component = () => {
    let updateCanceler: number;

    const backendApi = container.resolve(BackendApi);
    const [ipv4Leases, setIpv4Leases] = createSignal<Ipv4Lease[]>([]);
    const [newIpv4Leases, setNewIpv4Leases] = createSignal<Ipv4Lease[]>([]);
    const [removedIpv4Leases, setRemovedIpv4Leases] = createSignal<Ipv4Lease[]>([]);
    const [loading, setLoading] = createSignal(true);

    const fetchLeases = async () => {
        setLoading(true);

        const leasesResult = await backendApi.getIpv4Leases();

        setIpv4Leases(leasesResult.leases);
        setNewIpv4Leases(leasesResult.newLeases);
        setRemovedIpv4Leases(leasesResult.removedLeases);

        setLoading(false);
    };

    onMount(async () => {
        await fetchLeases();
        updateCanceler = setInterval(async () => {
            await fetchLeases();
        }, 5000);
    });

    onCleanup(() => {
        clearInterval(updateCanceler);
    });

    return (
        <>
            <span class="flex flex-row justify-center mb-8">IPv4 Leases</span>
            <div class={styles.loading_indicator}>
                {loading() && <span>Loading...</span>}
            </div>
            <div class={styles.lease_container}>
                <div class={styles.header}>
                    <span>IP Address</span>
                    <span>MAC Address</span>
                    <span>Expires At</span>
                </div>
                <div class={styles.newLeases}>
                    <IpLeaseBlock leases={newIpv4Leases()} type="new" />
                </div>
                <div id="currentLeases">
                    <IpLeaseBlock leases={ipv4Leases()} type="current" />
                </div>
                <div class={styles.removedLeases}>
                    <IpLeaseBlock leases={removedIpv4Leases()} type="removed" />
                </div>
            </div>
        </>
    );
};

export default App;
