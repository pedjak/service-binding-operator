Feature: Ready Service Binding Should Be Read-Only

    As a user of Service Binding operator
    I should not be able to modify a Service Binding that is ready

    Background:
        Given Namespace [TEST_NAMESPACE] is used
        * Service Binding Operator is running
        * CustomResourceDefinition backends.stable.example.com is available

    Scenario: Update a Service Binding that is ready
        Given The Custom Resource is present
        """
        apiVersion: stable.example.com/v1
        kind: Backend
        metadata:
            name: service-immutable
            annotations:
                service.binding/host: path={.spec.host}
        spec:
            host: foo
        """
        And Generic test application "app-immutable" is running
        And Service Binding is applied
            """
            apiVersion: binding.operators.coreos.com/v1alpha1
            kind: ServiceBinding
            metadata:
                name: sbr-immutable
            spec:
                services:
                  - group: stable.example.com
                    version: v1
                    kind: Backend
                    name: service-immutable
                application:
                    name: app-immutable
                    group: apps
                    version: v1
                    resource: deployments
            """
        When Service Binding "sbr-immutable" is ready
        Then Immutable Service Binding is unable to be applied
            """
            apiVersion: binding.operators.coreos.com/v1alpha1
            kind: ServiceBinding
            metadata:
                name: sbr-immutable
            spec:
                application:
                    name: app-immutable-2
                    group: apps
                    version: v1
                    resource: deployments
                services:
                  - group: stable.example.com
                    version: v1
                    kind: Backend
                    name: service-immutable
            """
