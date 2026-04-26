# Workout Service

[![release](https://img.shields.io/github/v/release/qkitzero/workout-service?logo=github)](https://github.com/qkitzero/workout-service/releases)
[![test](https://github.com/qkitzero/workout-service/actions/workflows/test.yml/badge.svg)](https://github.com/qkitzero/workout-service/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/qkitzero/workout-service/graph/badge.svg)](https://codecov.io/gh/qkitzero/workout-service)
[![Buf CI](https://github.com/qkitzero/workout-service/actions/workflows/buf-ci.yaml/badge.svg)](https://github.com/qkitzero/workout-service/actions/workflows/buf-ci.yaml)

- Microservices Architecture
- gRPC
- gRPC Gateway
- Buf ([buf.build/qkitzero-org/workout-service](https://buf.build/qkitzero-org/workout-service))
- Clean Architecture
- Docker
- Test
- Codecov
- Cloud Build
- Cloud Run

```mermaid
classDiagram
    direction LR

    class Set {
        id
        userID
        exerciseID
        rep
        weight
        trainedAt
        createdAt
    }

    class Exercise {
        id
        code
        category
        name
        muscles
    }

    class Muscle {
        id
        code
        name
    }

    Set "*" -- "1" Exercise
    Exercise "*" -- "*" Muscle
```

```mermaid
flowchart TD
    subgraph gcp[GCP]
        secret_manager[Secret Manager]

        subgraph cloud_build[Cloud Build]
            build_workout_service(Build workout-service)
            push_workout_service(Push workout-service)
            deploy_workout_service(Deploy workout-service)

            build_workout_service_gateway(Build workout-service-gateway)
            push_workout_service_gateway(Push workout-service-gateway)
            deploy_workout_service_gateway(Deploy workout-service-gateway)
        end


        subgraph artifact_registry[Artifact Registry]
            workout_service_image[(workout-service image)]
            workout_service_gateway_image[(workout-service-gateway image)]
        end

        subgraph cloud_run[Cloud Run]
            workout_service(Workout Service)
            workout_service_gateway(Workout Service Gateway)
        end
    end

    subgraph external[External]
        auth_service(Auth Service)
        workout_db[(Workout DB)]
    end

    build_workout_service --> push_workout_service --> workout_service_image
    build_workout_service_gateway --> push_workout_service_gateway --> workout_service_gateway_image

    workout_service_image --> deploy_workout_service --> workout_service
    workout_service_gateway_image --> deploy_workout_service_gateway --> workout_service_gateway

    secret_manager --> deploy_workout_service
    secret_manager --> deploy_workout_service_gateway

    workout_service_gateway --> workout_service
    workout_service --> workout_db
    workout_service --> auth_service
```
