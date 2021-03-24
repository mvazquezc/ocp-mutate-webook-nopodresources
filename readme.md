# OCP Mutation Webhook Test

This is a test repository.

The work here is heavily based on the work done by Alex Leonhardt here https://github.com/alex-leonhardt/k8s-mutate-webhook

More info here: https://medium.com/ovni/writing-a-very-basic-kubernetes-mutating-admission-webhook-398dbbcb63ec#:~:text=Mutating%20admission%20webhooks%20allow%20you,some%20some%20security%20requirements%2C%20etc.

This can be useful for testing your jsonpatches: https://json-schema-validator.herokuapp.com/jsonpatch.jsp

## Deploy the MutationWebhook on OpenShift

1. Deploy the webhook service

    ~~~sh
    oc create -f deploy/webhook-svc-deployment.yaml
    ~~~
2. Update the `CA_BUNDLE` for the webhook

    ~~~sh
    deploy/updatecabundle.sh deploy/webhook.yaml
    ~~~
3. Deploy the `MutatingWebhookConfiguration`

    ~~~sh
    oc create -f deploy/webhook.yaml
    ~~~
4. Deploy the test `namespace` and `deployment`

    > **NOTE**: As you can see the deployment has some requests and limits set, our mutator webhook will remove those.

    ~~~sh
    oc create -f deploy/test-app-deployment.yaml
    ~~~
