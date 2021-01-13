from kubernetes import client, config

# Configs can be set in Configuration class directly or using helper utility
config.load_kube_config()

v1 = client.CoreV1Api()
print("Listing pods with their IPs:")
# ret = v1.list_pod_for_all_namespaces(watch=False)
print("Pod %s does not exist. Creating it..." % name)
pod_manifest = {
    'apiVersion': 'v1',
    'kind': 'Pod',
    'metadata': {
        'name': name
    },
    'spec': {
        'containers': [{
            'image': 'busybox',
            'name': 'sleep',
            "args": [
                "/bin/sh",
                "-c",
                "while true;do date;sleep 5; done"
            ]
        }]
    }
}
resp = v1.create_namespaced_pod(body=pod_manifest,
                                          namespace='default')
for i in resp.items:
    print("%s\t%s\t%s" % (i.status.pod_ip, i.metadata.namespace, i.metadata.name))