#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>
#include <linux/pkt_cls.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/in.h>
#include <linux/tcp.h>


struct bpf_map_def SEC("maps") ip_tag_map = {
    .type = BPF_MAP_TYPE_HASH,
    .key_size = sizeof(__u32),
    .value_size = sizeof(__u32),
    .max_entries = 1024,
};

SEC("tcx")
int classify_tcp(struct __sk_buff *skb) {
    // Parse Ethernet header
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;

    // Parse Ethernet header
    struct ethhdr *eth = data;

    if ((void *)(eth + 1) > data_end) {
        return TC_ACT_OK; // Packet is incomplete
    }

    // Check if it's an IP packet
    if (eth->h_proto != __constant_htons(ETH_P_IP))
        return TC_ACT_OK;

    // Parse IP header
    struct iphdr *ip = (void *)(eth + 1);
    if ((void *)(ip + 1) > data_end) {
        return TC_ACT_OK; // Packet is incomplete
    }

    // Check if it's a TCP packet
    if (ip->protocol != IPPROTO_TCP)
        return TC_ACT_OK;

    __u32 src_ip = bpf_ntohs(ip->saddr);
    __u32 dest_ip = bpf_ntohs(ip->daddr);

    // Lookup source IP in the map
    __u32 *src_value = bpf_map_lookup_elem(&ip_tag_map, &src_ip);
    if (!src_value) {
        return TC_ACT_OK; // Source IP not found
    }

    // Lookup destination IP in the map
    __u32 *dest_value = bpf_map_lookup_elem(&ip_tag_map, &dest_ip);
    if (!dest_value) {
        return TC_ACT_OK; // Destination IP not found
    }

    int abs_value = (int) *src_value - (int) *dest_value;
    if (abs_value < 0) {
        abs_value = -abs_value;
    }
    if (abs_value > 0) {
        skb->tc_classid = abs_value + 10000; // Add magic number to avoid class overwrite
        bpf_printk("%pI4 %pI4", &ip->saddr, &ip->daddr);
    }


    return TC_ACT_OK; // Pass other packets without modification
}

char LICENSE[] SEC("license") = "GPL";