We have written a very simple Helm Chart to help you understand what you can and can't do with argocd-vault-replacer without breaking something in production.

You will need to deploy our third party Helm chart into your cluster using argocd (via git):

- Clone/copy the 'example' directory into your own git repo.
- Modify values.yaml so that the vault paths and secrets are valid for your environment.

Create a new Argo CD application to point to your repo:


You will need to change:
- spec.source.repoURL
- spec.source.path

```YAML
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: argocd-vault-replacer-example
spec:
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: example
  project: default
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
  source:
    repoURL: 'https://github.com/crumbhole/argocd-vault-replacer/'
    path: examples/example-third-party-helm-chart/example
    targetRevision: HEAD
    plugin:
      name: helm-argocd-vault-replacer
```

Apply your application to Argo CD in the usual way.

Argocd should then show a successful installation of the Helm chart into your cluster. You can then go and view the deployed artefacts and check whether the Vault secrets are as you expect them to be.
