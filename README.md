# Stubbies 🩳

Ethereum execution client stub for consensus layer clients.

## Getting Started

### Download a release
Download the latest release from the [Releases page](https://github.com/ethpandaops/stubbies/releases). Extract and run with:
```
./stubbies --config your-config.yaml
```

### Docker
Available as a docker image at [ethpandaops/stubbies](https://hub.docker.com/r/ethpandaops/stubbies/tags)

#### Images
- `latest` - distroless, multiarch
- `latest-debian` - debian, multiarch
- `$version` - distroless, multiarch, pinned to a release (i.e. `0.4.0`)
- `$version-debian` - debian, multiarch, pinned to a release (i.e. `0.4.0-debian`)

### Kubernetes via Helm
[Read more](https://github.com/ethpandaops/ethereum-helm-charts/tree/master/charts/stubbies)
```
helm repo add ethereum-helm-charts https://ethpandaops.github.io/ethereum-helm-charts

helm install stubbies ethereum-helm-charts/stubbies -f your_values.yaml
```

## Contact

Andrew - [@savid](https://twitter.com/Savid)

Sam - [@samcmau](https://twitter.com/samcmau)

