import { Component, For } from "solid-js"
import styles from './IpLeaseBlock.module.scss';

interface IpLeaseBlockProps {
    leases: Ipv4Lease[],
    type: 'current' | 'new' | 'removed'
}

function typeToLineClass(type: 'current' | 'new' | 'removed') {
    switch (type) {
        case 'current':
            return styles.current_lease_line;
        case 'new':
            return styles.new_lease_line;
        case 'removed':
            return styles.removed_lease_line;
    }
}

export const IpLeaseBlock: Component<IpLeaseBlockProps> = (props) => {
    return <>
        <For each={props.leases}>
            {lease => (
                <div class={typeToLineClass(props.type) + " mg-2"}>
                    <span>{lease.ipAddress}</span>
                    <span>{lease.hwAddress}</span>
                    <span>{new Date(lease.expiresAt * 1000).toLocaleString()}</span>
                </div>
            )}
        </For>
    </>
}
