#/bin/bash

k="docker exec -ti wavekit-k0s k0s kubectl"

echo "⏱️ waiting for k0s cluster to be ready..."
until $k get nodes | grep -q "Ready"; do
  sleep 1
done
echo "✅ k0s cluster is ready!"
echo ""
echo "Run the following command to get access to the cluster:"
echo ""
echo "   alias k='docker exec -ti wavekit-k0s k0s kubectl'"
echo ""
echo "The following commands are available:"
echo "🔍 make log         show the wavekit logs"
echo "🚀 make apply       apply the wavekit manifest in dev/k0s.yaml"
echo "🤞 make reload      reload the code when it changes"
echo ""
