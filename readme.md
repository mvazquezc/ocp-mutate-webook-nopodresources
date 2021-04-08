# OCP Mutation Webhook Test

This is a test repository.

The work here is heavily based on the work done by Alex Leonhardt here https://github.com/alex-leonhardt/k8s-mutate-webhook

More info here: https://medium.com/ovni/writing-a-very-basic-kubernetes-mutating-admission-webhook-398dbbcb63ec#:~:text=Mutating%20admission%20webhooks%20allow%20you,some%20some%20security%20requirements%2C%20etc.

This can be useful for testing your jsonpatches: https://json-schema-validator.herokuapp.com/jsonpatch.jsp

Useful info: https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/

## Deploy the MutationWebhook on OpenShift

1. Deploy the webhook service

    ~~~sh
    oc create -f deploy/webhook-svc-deployment.yaml
    ~~~
2. Update the `CA_BUNDLE` for the webhooks

    ~~~sh
    deploy/updatecabundle.sh deploy/mutatingwebhook.yaml
    deploy/updatecabundle.sh deploy/validatingwebhook.yaml
    ~~~
3. Deploy the `MutatingWebhookConfiguration` and the `ValidatingWebhookConfiguration`:

    ~~~sh
    oc create -f deploy/mutatingwebhook.yaml -f deploy/validatingwebhook.yaml
    ~~~
4. Test the webhooks:

    1. Mutation Webhook: 

        > **NOTE**: As you can see the deployment have requests and limits set to same values, our mutator webhook will do nothing.

        ~~~sh
        oc create -f deploy/test-mutating/test-app-deployment.yaml
        ~~~

        > **NOTE**: If the pod is Burstable (different requests and limits), our mutator webhook will request the removal of resources from pod.

        ~~~sh
        oc create -f deploy/test-mutating/test-app-deployment-burstable.yaml
        ~~~

        > **NOTE**: If the pod doesn't have any requests then no patch will be done

        ~~~sh
        oc create -f deploy/test-mutating/test-app-deployment-besteffort.yaml
        ~~~
    2. Validation Webhook

        > **NOTE**: As you can see the deployment have requests and limits set to same values, our validator will accept this request.

        ~~~sh
        oc create -f deploy/test-validating/test-app-deployment.yaml
        ~~~

        > **NOTE**: If the pod is Burstable (different requests and limits), our validator webhook will deny this request.

        ~~~sh
        oc create -f deploy/test-validating/test-app-deployment-burstable.yaml
        ~~~

        > **NOTE**: If the pod doesn't have any requests then the validator webhook will deny the request.

        ~~~sh
        oc create -f deploy/test-validating/test-app-deployment-besteffort.yaml
        ~~~

5. Clean everything:

    ~~~sh
    oc delete ns test-ns-mutate test-ns-validate
    ~~~

    ~~~sh
    oc delete -f deploy/webhook-svc-deployment.yaml 
    ~~~

    ~~~sh
    oc delete -f deploy/mutatingwebhook.yaml -f deploy/validatingwebhook.yaml
    ~~~
