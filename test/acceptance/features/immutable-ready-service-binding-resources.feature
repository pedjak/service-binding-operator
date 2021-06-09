Feature: Successful Service Binding Resources are Immutable

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

    Scenario: Allow modifying a Service Binding when it is not yet ready due to "Application Not Found"
        Given The Custom Resource is present
        """
        apiVersion: stable.example.com/v1
        kind: Backend
        metadata:
            name: service1
            annotations:
                service.binding/host: path={.spec.host}
        spec:
            host: foo
        """
        And Service Binding is applied
            """
            apiVersion: binding.operators.coreos.com/v1alpha1
            kind: ServiceBinding
            metadata:
                name: sbr-1
            spec:
                services:
                  - group: stable.example.com
                    version: v1
                    kind: Backend
                    name: service1
                application:
                    name: app1
                    group: apps
                    version: v1
                    resource: deployments
            """
        And jq ".status.conditions[] | select(.type=="Ready").status" of Service Binding "sbr-1" should be changed to "False"
        And The application "app1" does not exist
        When Generic test application "app2" is running
        And Service Binding is applied
            """
            apiVersion: binding.operators.coreos.com/v1alpha1
            kind: ServiceBinding
            metadata:
                name: sbr-1
            spec:
                services:
                  - group: stable.example.com
                    version: v1
                    kind: Backend
                    name: service1
                application:
                    name: app2
                    group: apps
                    version: v1
                    resource: deployments
            """
        Then Service Binding "sbr-1" is ready


    Scenario: Allow modifying a Service Binding when it is not yet ready due to "Service Not Found"
        Given Generic test application "app3" is running
        And The service "service2" does not exist
        And Service Binding is applied
            """
            apiVersion: binding.operators.coreos.com/v1alpha1
            kind: ServiceBinding
            metadata:
                name: sbr-2
            spec:
                services:
                  - group: stable.example.com
                    version: v1
                    kind: Backend
                    name: service2
                application:
                    name: app3
                    group: apps
                    version: v1
                    resource: deployments
            """
        And jq ".status.conditions[] | select(.type=="Ready").status" of Service Binding "sbr-2" should be changed to "False"
        And The Custom Resource is present
        """
        apiVersion: stable.example.com/v1
        kind: Backend
        metadata:
            name: service3
            annotations:
                service.binding/host: path={.spec.host}
        spec:
            host: foo
        """
        When Service Binding is applied
            """
            apiVersion: binding.operators.coreos.com/v1alpha1
            kind: ServiceBinding
            metadata:
                name: sbr-2
            spec:
                services:
                  - group: stable.example.com
                    version: v1
                    kind: Backend
                    name: service3
                application:
                    name: app3
                    group: apps
                    version: v1
                    resource: deployments
            """
        Then Service Binding "sbr-2" is ready
