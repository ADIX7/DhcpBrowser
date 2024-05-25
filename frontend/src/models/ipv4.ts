interface Ipv4LeasesResult {
    leases: Ipv4Lease[];
    newLeases: Ipv4Lease[];
    removedLeases: Ipv4Lease[];
}

interface Ipv4Lease {
    ipAddress: string;
    hwAddress: string;
    expiresAt: number;
}
