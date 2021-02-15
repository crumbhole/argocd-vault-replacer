# vault-replacer

Mainly intended as a plugin for argocd

Will scan the current directory recursively for any .yaml (or .yml if you're so inclined) files and attempt to replaces strings of the form <vault:/store/data/path!key> with those obtained from a vault kv2 store. It is intended that this is run from within your argocd-server pod as a plugin.

## Authentication

It only has two methods of authenticating with vault:
* Using kubernetes authentication method https://github.com/hashicorp/vault/blob/master/website/content/docs/auth/kubernetes.mdx
* Using a token, which is only intended for debugging

Both methods expect the environment variable VAULT_ADDR to be set.

It will attempt to use kubernetes authentication through an appropriate service account first, and complain if that doesn't work. It will then use VAULT_TOKEN which should be a valid token. This tool has no way of renewing a token or obtaining one other than through a kubernetes service account.

To use the kubernetes service account your pod should be running with the appropriate service account, and will try to obtain the JWT token from /var/run/secrets/kubernetes.io/serviceaccount/token which is the default location.

It will use the environment variable VAULT_ROLE as the name of the role for that token, defaulting to "argocd".
It will use the environment variable VAULT_AUTH_PATH to determine the authorisation path for kubernetes authentication. This defaults in this tool and in vault to "kubernetes" so will probably not need configuring.

## Valid authentication paths

Currently the only valid 'URL style' to a path is

<vault:/store/data/path!key>

You must put the ../data/.. into the path. If your path or key contains !, <, > or | you must URL escape it. If you path or key has one or more leading or trailing spaces or tabs you must URL escape them you weirdo.

